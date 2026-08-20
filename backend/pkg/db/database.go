package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database connection parameters.
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

var dbConn *gorm.DB

// InitDB establishes a MySQL connection using GORM.
func InitDB(cfg Config) (*gorm.DB, error) {
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	maxIdle := cfg.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 10
	}
	maxOpen := cfg.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 100
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbConn = db
	return db, nil
}

// GetDB returns the database connection.
func GetDB() *gorm.DB {
	return dbConn
}

// AutoMigrate runs automatic migration for all registered models.
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// systemSetting is used to read system settings from DB.
type systemSetting struct {
	Key   string `gorm:"column:key"`
	Value string `gorm:"column:value"`
	Group string `gorm:"column:group"`
	Name  string `gorm:"column:name"`
}

func (systemSetting) TableName() string {
	return "system_configs"
}

// GetSystemSetting reads a single setting value from the database.
func GetSystemSetting(key string) string {
	if dbConn == nil {
		return ""
	}
	var s systemSetting
	if err := dbConn.Where("`key` = ?", key).First(&s).Error; err != nil {
		return ""
	}
	return s.Value
}

// GetSystemSettings reads all settings from a group.
func GetSystemSettings(group string) map[string]string {
	if dbConn == nil {
		return nil
	}
	var settings []systemSetting
	if err := dbConn.Where("`group` = ?", group).Find(&settings).Error; err != nil {
		return nil
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}

// SetSystemSetting updates or creates a setting.
func SetSystemSetting(key, value, group, desc string) error {
	if dbConn == nil {
		return fmt.Errorf("database not initialized")
	}
	s := systemSetting{Key: key, Value: value, Group: group, Name: desc}
	result := dbConn.Where("`key` = ?", key).Assign(map[string]interface{}{
		"value":  value,
		"group":  group,
		"name":   desc,
	}).FirstOrCreate(&s)
	return result.Error
}
