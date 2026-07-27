package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/anchor-finance/backend/internal/api/middleware"
	"github.com/anchor-finance/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Server holds the application server dependencies.
type Server struct {
	db          *gorm.DB
	redis       *redis.Client
	jwtManager  *auth.JWTManager
	router      *gin.Engine
	installH    *InstallHandler
}

// NewServer creates and configures the HTTP server with all routes and middleware.
func NewServer(dbConn *gorm.DB, rdb *redis.Client, jwtSecret string) *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		db:         dbConn,
		redis:      rdb,
		jwtManager: auth.NewJWTManager(jwtSecret, 72),
		router:     gin.New(),
		installH:   NewInstallHandler(),
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

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true, "message": "pong"})
		})

		// Public routes (no auth required)
		public := v1.Group("")
		{
			public.POST("/login", s.handleLogin)
			public.POST("/register", s.handleRegister)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(s.jwtManager))
		{
			protected.GET("/user/profile", s.handleGetProfile)
			protected.PUT("/user/profile", s.handleUpdateProfile)
		}
	}

	// API v2 (reserved for future use)
	v2 := s.router.Group("/api/v2")
	{
		v2.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true, "message": "pong v2"})
		})
	}

	// Admin routes
	admin := s.router.Group("/api/admin")
	admin.Use(middleware.JWTAuth(s.jwtManager))
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/users", s.handleAdminListUsers)
		admin.GET("/settings", s.handleAdminGetSettings)
		admin.PUT("/settings", s.handleAdminUpdateSettings)
		admin.GET("/dashboard", s.handleAdminDashboard)
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

// --------------- placeholder handlers ---------------

func (s *Server) handleLogin(c *gin.Context) {
	// TODO: implement login logic
	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "login endpoint"})
}

func (s *Server) handleRegister(c *gin.Context) {
	// TODO: check if registration is enabled
	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "register endpoint"})
}

func (s *Server) handleGetProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextKeyUserID)
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"user_id": userID,
		"message": "profile endpoint",
	})
}

func (s *Server) handleUpdateProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextKeyUserID)
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"user_id": userID,
		"message": "update profile endpoint",
	})
}

func (s *Server) handleAdminListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "admin list users endpoint",
	})
}

func (s *Server) handleAdminGetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "admin get settings endpoint",
	})
}

func (s *Server) handleAdminUpdateSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "admin update settings endpoint",
	})
}

func (s *Server) handleAdminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "admin dashboard endpoint",
	})
}
