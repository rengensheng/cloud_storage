package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"cloud-storage/internal/config"
	"cloud-storage/internal/database"
	"cloud-storage/internal/handlers"
	"cloud-storage/internal/middleware"
	"cloud-storage/internal/pkg/storage"
	"cloud-storage/internal/repositories"
	"cloud-storage/internal/services"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置日志
	setupLogging(cfg)

	// 初始化数据库
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// 初始化Redis
	_, err = database.InitRedis(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
		log.Println("Continuing without Redis support")
	} else {
		defer database.CloseRedis()
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		log.Printf("Warning: Failed to migrate database: %v", err)
	}

	// 初始化存储
	storageImpl, err := setupStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// 初始化仓库
	fileRepo := repositories.NewFileRepository(db)
	userRepo := repositories.NewUserRepository(db)
	shareRepo := repositories.NewShareRepository(db)
	operationLogRepo := repositories.NewOperationLogRepository(db)

	// 初始化服务
	fileService := services.NewFileService(cfg, db, fileRepo, userRepo, storageImpl)
	shareService := services.NewShareService(db, shareRepo, fileRepo)
	operationLogService := services.NewOperationLogService(operationLogRepo)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// 初始化处理器
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(&userRepo, authMiddleware)
	shareHandler := handlers.NewShareHandler(shareService)
	adminHandler := handlers.NewAdminHandler(userRepo, operationLogService, shareService, fileService)

	// 设置Gin模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由器
	router := gin.New()

	// 注册中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		// 公开路由
		public := api.Group("")
		authHandler.RegisterRoutes(public)

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(authMiddleware.Authenticate())
		fileHandler.RegisterRoutes(protected)
		shareHandler.RegisterRoutes(protected, public)
		adminHandler.RegisterRoutes(protected)
	}

	// 启动服务器
	startServer(cfg, router)
}

// setupLogging 设置日志
func setupLogging(cfg *config.Config) {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Warning: Failed to create logs directory: %v", err)
	}

	// 设置日志输出
	if cfg.Log.File != "" {
		logFile, err := os.OpenFile(cfg.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Warning: Failed to open log file: %v", err)
		} else {
			log.SetOutput(logFile)
		}
	}

	log.Printf("Starting cloud storage service in %s mode", cfg.App.Env)
}

// setupStorage 设置存储
func setupStorage(cfg *config.Config) (storage.Storage, error) {
	storageConfig := storage.StorageConfig{
		Type:      storage.StorageTypeLocal,
		LocalPath: cfg.Storage.StoragePath,
	}

	// 创建存储实例
	storageImpl, err := storage.NewStorage(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	// 创建必要的目录
	if err := os.MkdirAll(cfg.Storage.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	if err := os.MkdirAll(cfg.Storage.TempPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	log.Printf("Storage initialized at: %s", cfg.Storage.StoragePath)
	return storageImpl, nil
}

// startServer 启动服务器
func startServer(cfg *config.Config, router *gin.Engine) {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// init 初始化函数
func init() {
	// 设置时区
	os.Setenv("TZ", "Asia/Shanghai")
}
