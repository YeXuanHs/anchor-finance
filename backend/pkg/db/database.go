package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database connection parameters.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var dbConn *gorm.DB

// InitDB establishes a PostgreSQL connection using GORM.
func InitDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
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
}

func (systemSetting) TableName() string {
	return "system_settings"
}

// GetSystemSetting reads a single setting value from the database.
func GetSystemSetting(key string) string {
	if dbConn == nil {
		return ""
	}
	var s systemSetting
	if err := dbConn.Where("key = ?", key).First(&s).Error; err != nil {
		return ""
	}
	return s.Value
}

// SetSystemSetting writes a single setting value to the database.
func SetSystemSetting(key, value string) error {
	if dbConn == nil {
		return fmt.Errorf("database not initialized")
	}
	return dbConn.Where("key = ?", key).Assign(systemSetting{Value: value}).FirstOrCreate(&systemSetting{Key: key}).Error
}
