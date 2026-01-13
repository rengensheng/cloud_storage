package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"cloud-storage/internal/config"
	"cloud-storage/internal/models"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.Timezone,
	)

	// 配置GORM日志
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Database connection established successfully")
	return db, nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database connection: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// 导入模型包以进行迁移
	importModels()

	log.Println("Starting database migration...")

	// 执行迁移
	err := DB.AutoMigrate(
		// 用户相关
		&models.User{},
		&models.LoginAttempt{},

		// 文件相关
		&models.File{},
		&models.FileVersion{},

		// 分享相关
		&models.Share{},

		// 日志相关
		&models.OperationLog{},
		&models.SecurityAlert{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// importModels 导入模型包，确保模型被注册
func importModels() {
	// 这里只需要导入模型包，GORM会自动发现模型
	// 实际的模型将在models包中定义
}

// CreateDatabase 创建数据库（如果不存在）
func CreateDatabase(cfg *config.Config) error {
	// 连接到默认数据库以创建目标数据库
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to default database: %w", err)
	}

	// 检查数据库是否存在
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", cfg.Database.Name).Scan(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// 如果数据库不存在，则创建
	if count == 0 {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name)).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", cfg.Database.Name)
	} else {
		log.Printf("Database '%s' already exists", cfg.Database.Name)
	}

	// 关闭临时连接
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return nil
}