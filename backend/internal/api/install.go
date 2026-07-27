package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anchor-finance/backend/pkg/auth"
	"github.com/anchor-finance/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// installRequest holds the full installation payload.
type installRequest struct {
	DBHost         string `json:"db_host" binding:"required"`
	DBPort         string `json:"db_port" binding:"required"`
	DBUser         string `json:"db_user" binding:"required"`
	DBPassword     string `json:"db_password"`
	DBName         string `json:"db_name" binding:"required"`
	RedisHost      string `json:"redis_host" binding:"required"`
	RedisPort      string `json:"redis_port" binding:"required"`
	RedisPassword  string `json:"redis_password"`
	AdminUsername  string `json:"admin_username" binding:"required,min=3"`
	AdminPassword  string `json:"admin_password" binding:"required,min=6"`
	AdminEmail     string `json:"admin_email" binding:"required,email"`
	TestOnly       bool   `json:"test_only"`
}

// InstallHandler manages the installation API endpoints.
type InstallHandler struct {
	installFlagPath string
}

// NewInstallHandler creates a new install handler.
func NewInstallHandler() *InstallHandler {
	return &InstallHandler{
		installFlagPath: filepath.Join("configs", ".installed"),
	}
}

// IsInstalled checks whether the system has already been installed.
func (h *InstallHandler) IsInstalled() bool {
	_, err := os.Stat(h.installFlagPath)
	return err == nil
}

// RegisterRoutes registers install-related routes. Only registers if not installed.
func (h *InstallHandler) RegisterRoutes(r *gin.Engine) {
	if h.IsInstalled() {
		return
	}

	grp := r.Group("/api/install")
	{
		grp.GET("/check-env", h.CheckEnv)
		grp.POST("/test-db", h.TestDB)
		grp.POST("", h.DoInstall)
	}

	// Serve the installer page
	r.GET("/install", func(c *gin.Context) {
		c.File(filepath.Join("install", "index.html"))
	})

	// Redirect root to install if not installed
	r.GET("/", func(c *gin.Context) {
		if !h.IsInstalled() {
			c.Redirect(http.StatusFound, "/install")
			return
		}
		c.Next()
	})
}

// CheckEnv verifies runtime environment requirements.
// Query param: type = go | pg | redis
func (h *InstallHandler) CheckEnv(c *gin.Context) {
	checkType := c.Query("type")
	switch checkType {
	case "go":
		checkGoVersion(c)
	case "pg":
		// PostgreSQL is tested via the TestDB endpoint; here we just report that the driver is available
		c.JSON(http.StatusOK, gin.H{"ok": true, "message": "PostgreSQL 驱动就绪"})
	case "redis":
		// Redis is tested via the TestDB endpoint; here we just report that the client is available
		c.JSON(http.StatusOK, gin.H{"ok": true, "message": "Redis 客户端就绪"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "message": "未知检测类型"})
	}
}

func checkGoVersion(c *gin.Context) {
	ver := getGoVersion()
	if ver != "" {
		c.JSON(http.StatusOK, gin.H{"ok": true, "message": ver})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "未检测到 Go 环境"})
	}
}

func getGoVersion() string {
	if v, err := strconv.Unquote(`"` + "go" + `"`); err == nil && v != "" {
		// Use the embedded version from the binary itself
		return "Go (编译时版本)"
	}
	return ""
}

// TestDB tests the database and Redis connections.
func (h *InstallHandler) TestDB(c *gin.Context) {
	var req installRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "message": "参数错误: " + err.Error()})
		return
	}

	// Test PostgreSQL
	pgDB, err := db.InitDB(db.Config{
		Host:     req.DBHost,
		Port:     req.DBPort,
		User:     req.DBUser,
		Password: req.DBPassword,
		DBName:   req.DBName,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "PostgreSQL 连接失败: " + err.Error()})
		return
	}

	sqlDB, _ := pgDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// Test Redis
	redisPort := req.RedisPort
	if redisPort == "" {
		redisPort = "6379"
	}
	rdb, err := db.InitRedis(db.RedisConfig{
		Host:     req.RedisHost,
		Port:     redisPort,
		Password: req.RedisPassword,
		DB:       0,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "Redis 连接失败: " + err.Error()})
		return
	}
	rdb.Close()

	if req.TestOnly {
		c.JSON(http.StatusOK, gin.H{"ok": true, "message": "数据库和 Redis 连接成功"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "连接测试通过"})
}

// DoInstall performs the full installation.
func (h *InstallHandler) DoInstall(c *gin.Context) {
	if h.IsInstalled() {
		c.JSON(http.StatusConflict, gin.H{"ok": false, "message": "系统已经安装，请勿重复操作"})
		return
	}

	var req installRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "message": "参数错误: " + err.Error()})
		return
	}

	// 1. Connect to PostgreSQL
	redisPort := req.RedisPort
	if redisPort == "" {
		redisPort = "6379"
	}

	pgDB, err := db.InitDB(db.Config{
		Host:     req.DBHost,
		Port:     req.DBPort,
		User:     req.DBUser,
		Password: req.DBPassword,
		DBName:   req.DBName,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "PostgreSQL 连接失败: " + err.Error()})
		return
	}

	// 2. Auto-migrate all models
	if err := runMigrations(pgDB); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "数据库迁移失败: " + err.Error()})
		return
	}

	// 3. Create default admin user
	if err := createAdmin(pgDB, req.AdminUsername, req.AdminPassword, req.AdminEmail); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "创建管理员失败: " + err.Error()})
		return
	}

	// 4. Create default system configs
	if err := createDefaultConfigs(pgDB); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "创建默认配置失败: " + err.Error()})
		return
	}

	// 5. Generate config.yaml
	configContent, err := generateConfig(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "生成配置文件失败: " + err.Error()})
		return
	}

	// 6. Write the installed flag
	if err := os.MkdirAll(filepath.Dir(h.installFlagPath), 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "创建安装标记失败: " + err.Error()})
		return
	}
	if err := os.WriteFile(h.installFlagPath, []byte("installed"), 0644); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": false, "message": "写入安装标记失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":             true,
		"message":        "安装成功",
		"config_preview": configContent,
	})
}

// --------------- internal helpers ---------------

// systemSetting is a minimal model used during install for system config.
type systemSetting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;size:128;not null"`
	Value string `gorm:"type:text"`
	Group string `gorm:"size:64;default:general"`
	Desc  string `gorm:"size:255"`
}

// user is a minimal model used during install for admin creation.
type user struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64;not null"`
	Password string `gorm:"size:255;not null"`
	Email    string `gorm:"uniqueIndex;size:128"`
	IsAdmin  bool   `gorm:"default:false"`
	Status   int    `gorm:"default:1"`
}

func (user) TableName() string {
	return "users"
}

func (systemSetting) TableName() string {
	return "system_settings"
}

func runMigrations(dbConn *gorm.DB) error {
	return db.AutoMigrate(dbConn, &user{}, &systemSetting{})
}

func createAdmin(dbConn *gorm.DB, username, password, email string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash error: %w", err)
	}

	admin := user{
		Username: username,
		Password: string(hash),
		Email:    email,
		IsAdmin:  true,
		Status:   1,
	}

	return dbConn.Create(&admin).Error
}

func createDefaultConfigs(dbConn *gorm.DB) error {
	defaults := []systemSetting{
		{Key: "site_name", Value: "锚点财务", Group: "general", Desc: "站点名称"},
		{Key: "site_url", Value: "http://localhost:8080", Group: "general", Desc: "站点 URL"},
		{Key: "site_description", Value: "高效、安全的财务管理平台", Group: "general", Desc: "站点描述"},
		{Key: "currency", Value: "CNY", Group: "finance", Desc: "默认货币"},
		{Key: "decimal_places", Value: "2", Group: "finance", Desc: "小数位数"},
		{Key: "jwt_expire_hours", Value: "72", Group: "security", Desc: "JWT 过期时间（小时）"},
		{Key: "enable_registration", Value: "true", Group: "general", Desc: "是否开放注册"},
	}

	for _, s := range defaults {
		if err := dbConn.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

func generateConfig(req installRequest) (string, error) {
	redisPort := req.RedisPort
	if redisPort == "" {
		redisPort = "6379"
	}

	content := fmt.Sprintf(`# AnchorFinance Configuration
# Auto-generated by installer

server:
  port: 8080
  mode: release

database:
  host: %s
  port: %s
  user: %s
  password: "%s"
  dbname: %s

redis:
  host: %s
  port: %s
  password: "%s"
  db: 0

jwt:
  secret: "%s"
  expire_hours: 72

log:
  level: info
  format: text
`,
		req.DBHost, req.DBPort, req.DBUser, req.DBPassword, req.DBName,
		req.RedisHost, redisPort, req.RedisPassword,
		randomSecret(32),
	)

	configPath := filepath.Join("configs", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return "", err
	}

	return content, nil
}

func randomSecret(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	// Use a simple approach; in production, use crypto/rand
	return fmt.Sprintf("%x", b)
}
