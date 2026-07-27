package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/anchor-finance/backend/internal/api/admin"
	"github.com/anchor-finance/backend/internal/api/middleware"
	v1 "github.com/anchor-finance/backend/internal/api/v1"
	v2 "github.com/anchor-finance/backend/internal/api/v2"
	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/auth"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Deps holds all dependencies for route registration.
type Deps struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Log     *logger.Logger
	JWTKey  string
	UserSvc *service.UserService
	ProdSvc *service.ProductService
	OrdSvc  *service.OrderService
	InvSvc  *service.InvoiceService
	TicSvc  *service.TicketService
	PaySvc  *service.PaymentService
	CartSvc *service.CartService
	OAuthSvc *service.OAuthService
}

// Server holds the application server dependencies.
type Server struct {
	db         *gorm.DB
	redis      *redis.Client
	jwtManager *auth.JWTManager
	router     *gin.Engine
	installH   *InstallHandler
	deps       *Deps
}

// NewServer creates and configures the HTTP server with all routes and middleware.
func NewServer(dbConn *gorm.DB, rdb *redis.Client, jwtSecret string) *Server {
	gin.SetMode(gin.ReleaseMode)

	log := logger.New()
	jwtMgr := auth.NewJWTManager(jwtSecret, 72)

	// Initialize all services
	userSvc := service.NewUserService(dbConn, log)
	prodSvc := service.NewProductService(dbConn, log)
	ordSvc := service.NewOrderService(dbConn, log)
	invSvc := service.NewInvoiceService(dbConn, log)
	ticSvc := service.NewTicketService(dbConn, log)
	paySvc := service.NewPaymentService(dbConn, log)
	cartSvc := service.NewCartService(dbConn, log)
	oauthSvc := service.NewOAuthService(dbConn, log)

	deps := &Deps{
		DB:       dbConn,
		Redis:    rdb,
		Log:      log,
		JWTKey:   jwtSecret,
		UserSvc:  userSvc,
		ProdSvc:  prodSvc,
		OrdSvc:   ordSvc,
		InvSvc:   invSvc,
		TicSvc:   ticSvc,
		PaySvc:   paySvc,
		CartSvc:  cartSvc,
		OAuthSvc: oauthSvc,
	}

	s := &Server{
		db:         dbConn,
		redis:      rdb,
		jwtManager: jwtMgr,
		router:     gin.New(),
		installH:   NewInstallHandler(),
		deps:       deps,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware registers global middleware.
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.CORS())
}

// setupRoutes registers all route groups.
func (s *Server) setupRoutes() {
	// Install routes (only active when not installed)
	s.installH.RegisterRoutes(s.router)

	// Serve frontend static files
	s.serveStaticFiles()

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 - 兼容智简魔方
	v1Group := s.router.Group("/api/v1")
	v1.RegisterRoutes(v1Group, s.deps.toV1Deps())

	// API v2 - 锚点财务原生
	v2Group := s.router.Group("/api/v2")
	v2.RegisterRoutes(v2Group, s.deps.toV2Deps())

	// Admin routes
	adminGroup := s.router.Group("/api/admin")
	admin.RegisterRoutes(adminGroup, s.deps.toAdminDeps())
}

// toV1Deps converts Deps to v1.Deps.
func (d *Deps) toV1Deps() v1.Deps {
	return v1.Deps{
		DB:      d.DB,
		Log:     d.Log,
		JWTKey:  d.JWTKey,
		UserSvc: d.UserSvc,
		ProdSvc: d.ProdSvc,
		OrdSvc:  d.OrdSvc,
		InvSvc:  d.InvSvc,
		TicSvc:  d.TicSvc,
	}
}

// toV2Deps converts Deps to v2.Deps.
func (d *Deps) toV2Deps() v2.Deps {
	return v2.Deps{
		DB:      d.DB,
		Log:     d.Log,
		JWTKey:  d.JWTKey,
		UserSvc: d.UserSvc,
		ProdSvc: d.ProdSvc,
		OrdSvc:  d.OrdSvc,
		InvSvc:  d.InvSvc,
		TicSvc:  d.TicSvc,
		CartSvc: d.CartSvc,
	}
}

// toAdminDeps converts Deps to admin.Deps.
func (d *Deps) toAdminDeps() admin.Deps {
	return admin.Deps{
		DB:      d.DB,
		Log:     d.Log,
		JWTKey:  d.JWTKey,
		UserSvc: d.UserSvc,
		ProdSvc: d.ProdSvc,
		OrdSvc:  d.OrdSvc,
		InvSvc:  d.InvSvc,
		TicSvc:  d.TicSvc,
	}
}

// serveStaticFiles configures static file serving for frontend and admin panels.
func (s *Server) serveStaticFiles() {
	frontendDir := "frontend"
	adminDir := "admin"

	// Serve admin static files at /admin
	if info, err := os.Stat(adminDir); err == nil && info.IsDir() {
		s.router.Static("/admin", adminDir)
		s.router.NoRoute(func(c *gin.Context) {
			// If the path starts with /admin, serve admin index.html for SPA routing
			if len(c.Request.URL.Path) >= 6 && c.Request.URL.Path[:6] == "/admin" {
				indexPath := filepath.Join(adminDir, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					c.File(indexPath)
					return
				}
			}
			// For all other routes, serve frontend index.html (SPA fallback)
			indexPath := filepath.Join(frontendDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"ok": false, "message": "not found"})
		})
	}

	// Serve frontend static files at /
	if info, err := os.Stat(frontendDir); err == nil && info.IsDir() {
		s.router.Static("/assets", filepath.Join(frontendDir, "assets"))
		s.router.StaticFile("/favicon.ico", filepath.Join(frontendDir, "favicon.ico"))
	}
}

// Run starts the HTTP server on the given address.
func (s *Server) Run(addr string) error {
	logrus.Infof("AnchorFinance server starting on %s", addr)
	return s.router.Run(addr)
}

// Router returns the underlying gin.Engine for testing purposes.
func (s *Server) Router() *gin.Engine {
	return s.router
}
