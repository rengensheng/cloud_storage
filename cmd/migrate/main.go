package main

import (
	"flag"
	"log"
	"os"

	"cloud-storage/internal/config"
	"cloud-storage/internal/database"
	"cloud-storage/internal/models"
	"gorm.io/gorm"
)

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	var rollback bool
	flag.BoolVar(&rollback, "rollback", false, "rollback migrations")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// 执行迁移或回滚
	if rollback {
		rollbackMigrations(db)
	} else {
		runMigrations(db)
	}
}

// runMigrations 运行数据库迁移
func runMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.FileVersion{},
		&models.Share{},
		&models.OperationLog{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 创建索引
	if err := createIndexes(db); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	// 创建默认管理员账户
	if err := createDefaultAdmin(db); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	log.Println("Migrations completed successfully!")
}

// rollbackMigrations 回滚数据库迁移
func rollbackMigrations(db *gorm.DB) {
	log.Println("Rolling back database migrations...")
	log.Println("Warning: AutoMigrate doesn't support rollback, you need to manually drop tables")
	log.Println("Tables to drop:")
	log.Println("  - operation_logs")
	log.Println("  - shares")
	log.Println("  - file_versions")
	log.Println("  - files")
	log.Println("  - users")
	log.Println("Please use SQL commands to drop these tables manually.")
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) error {
	log.Println("Creating indexes...")

	// 用户表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)").Error; err != nil {
		return err
	}

	// 文件表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_user_id ON files(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_parent_id ON files(parent_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_user_id_parent_id_name ON files(user_id, parent_id, name)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_share_token ON files(share_token)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_type ON files(type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files(deleted_at)").Error; err != nil {
		return err
	}

	// 文件版本表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_file_versions_file_id ON file_versions(file_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_file_versions_version_number ON file_versions(file_id, version_number)").Error; err != nil {
		return err
	}

	// 分享表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_file_id ON shares(file_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_user_id ON shares(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_share_token ON shares(share_token)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_expires_at ON shares(expires_at)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shares_is_active ON shares(is_active)").Error; err != nil {
		return err
	}

	// 操作日志表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_operation_logs_operation_type ON operation_logs(operation_type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_operation_logs_resource_id ON operation_logs(resource_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at)").Error; err != nil {
		return err
	}

	log.Println("Indexes created successfully!")
	return nil
}

// createDefaultAdmin 创建默认管理员账户
func createDefaultAdmin(db *gorm.DB) error {
	log.Println("Creating default admin account...")

	// 检查是否已存在管理员
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin account already exists, skipping...")
		return nil
	}

	// 从环境变量读取管理员信息
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" {
		adminUsername = "admin"
	}
	if adminEmail == "" {
		adminEmail = "admin@cloud-storage.local"
	}
	if adminPassword == "" {
		adminPassword = "admin123456"
	}

	// 创建管理员
	admin := &models.User{
		Username:     adminUsername,
		Email:        adminEmail,
		PasswordHash: adminPassword,
		Role:         "admin",
		StorageQuota: 107374182400, // 100GB
		IsActive:     true,
	}

	// 哈希密码会通过模型的BeforeCreate钩子自动处理
	if err := db.Create(admin).Error; err != nil {
		return err
	}

	log.Printf("Default admin account created successfully!")
	log.Printf("Username: %s", adminUsername)
	log.Printf("Email: %s", adminEmail)
	log.Printf("Password: %s", adminPassword)
	log.Printf("Please change the default password after first login!")

	return nil
}
