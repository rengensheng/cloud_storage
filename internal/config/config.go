package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Security SecurityConfig
	Log      LogConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Env  string
	Name string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	Timezone string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret              string
	ExpireHours         int
	RefreshExpireHours  int
}

// StorageConfig 存储配置
type StorageConfig struct {
	StoragePath      string
	TempPath         string
	MaxUploadSize    int64
	MaxMemorySize    int64
	EnableChunkUpload bool
	ChunkSize        int64
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORSAllowOrigins   string
	CORSAllowCredentials bool
	RateLimit          int
	RateLimitDuration  time.Duration
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string
	File     string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Name: getEnv("APP_NAME", "cloud-storage"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "cloud_storage"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Shanghai"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
			ExpireHours:        getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpireHours: getEnvAsInt("JWT_REFRESH_EXPIRE_HOURS", 168),
		},
		Storage: StorageConfig{
			StoragePath:      getEnv("STORAGE_PATH", "./storage/uploads"),
			TempPath:         getEnv("TEMP_PATH", "./storage/temp"),
			MaxUploadSize:    getEnvAsInt64("MAX_UPLOAD_SIZE", 104857600),  // 100MB
			MaxMemorySize:    getEnvAsInt64("MAX_MEMORY_SIZE", 33554432),   // 32MB
			EnableChunkUpload: getEnvAsBool("ENABLE_CHUNK_UPLOAD", true),
			ChunkSize:        getEnvAsInt64("CHUNK_SIZE", 5242880),         // 5MB
		},
		Security: SecurityConfig{
			CORSAllowOrigins:   getEnv("CORS_ALLOW_ORIGINS", "*"),
			CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			RateLimit:          getEnvAsInt("RATE_LIMIT", 100),
			RateLimitDuration:  time.Duration(getEnvAsInt("RATE_LIMIT_DURATION", 60)) * time.Second,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "./logs/app.log"),
		},
	}
}

// 辅助函数：获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}