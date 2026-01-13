package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"cloud-storage/internal/config"
	"cloud-storage/internal/database"
)

// Claims JWT声明
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	cfg *config.Config
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// Authenticate 认证中间件
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer令牌格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析和验证JWT令牌
		claims, err := m.parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// 检查令牌是否在Redis黑名单中（如果支持注销）
		if m.isTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// ParseToken 解析JWT令牌
func (m *AuthMiddleware) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// parseToken 解析JWT令牌
func (m *AuthMiddleware) parseToken(tokenString string) (*Claims, error) {
	return m.ParseToken(tokenString)
}

// isTokenBlacklisted 检查令牌是否在黑名单中
func (m *AuthMiddleware) isTokenBlacklisted(tokenString string) bool {
	// 计算令牌的哈希作为键
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])
	key := fmt.Sprintf("blacklist:token:%s", tokenHash)

	// 检查Redis中是否存在
	exists, err := database.Exists(key)
	if err != nil {
		// Redis错误，默认认为令牌有效
		return false
	}

	return exists
}

// GenerateToken 生成JWT令牌
func (m *AuthMiddleware) GenerateToken(userID uuid.UUID, username, role string) (string, error) {
	// 设置令牌过期时间
	expireTime := time.Now().Add(time.Duration(m.cfg.JWT.ExpireHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cloud-storage",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.JWT.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func (m *AuthMiddleware) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	// 刷新令牌使用更长的过期时间
	expireTime := time.Now().Add(time.Duration(m.cfg.JWT.RefreshExpireHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cloud-storage",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.JWT.Secret))
}

// RefreshToken 刷新访问令牌
func (m *AuthMiddleware) RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// 生成新的访问令牌和刷新令牌
	newAccessToken, err := m.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := m.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// BlacklistToken 将令牌加入黑名单
func (m *AuthMiddleware) BlacklistToken(tokenString string, expireTime time.Time) error {
	// 计算令牌的哈希作为键
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])
	key := fmt.Sprintf("blacklist:token:%s", tokenHash)

	// 计算剩余过期时间
	expiration := time.Until(expireTime)
	if expiration <= 0 {
		// 令牌已过期，不需要加入黑名单
		return nil
	}

	// 将令牌哈希存储到Redis，设置与令牌相同的过期时间
	return database.Set(key, "1", expiration)
}

// RequireRole 要求特定角色的中间件
func (m *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// 检查用户角色
		if !hasPermission(userRole.(string), requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission 检查用户是否有权限
func hasPermission(userRole, requiredRole string) bool {
	// 角色权限层级
	roleHierarchy := map[string]int{
		"user":  1,
		"admin": 2,
	}

	userLevel, userOk := roleHierarchy[userRole]
	requiredLevel, requiredOk := roleHierarchy[requiredRole]

	if !userOk || !requiredOk {
		return false
	}

	return userLevel >= requiredLevel
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *gin.Context) (uuid.UUID, string, string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, "", "", fmt.Errorf("user not authenticated")
	}

	username, _ := c.Get("username")
	role, _ := c.Get("role")

	return userID.(uuid.UUID), username.(string), role.(string), nil
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer令牌格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		// 尝试解析令牌
		claims, err := m.parseToken(tokenString)
		if err != nil {
			// 令牌无效，继续处理（作为未认证用户）
			c.Next()
			return
		}

		// 检查令牌是否在黑名单中
		if m.isTokenBlacklisted(tokenString) {
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		// 构建Redis键
		key := fmt.Sprintf("ratelimit:%s:%s", c.FullPath(), clientIP)

		// 检查速率限制
		allowed, err := database.RateLimit(key, limit, window)
		if err != nil {
			// Redis错误，跳过速率限制
			c.Next()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置CORS头
		if cfg.Security.CORSAllowOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.Security.CORSAllowOrigins)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取用户ID（如果已认证）
		userID := "anonymous"
		if id, exists := c.Get("userID"); exists {
			userID = id.(uuid.UUID).String()
		}

		// 记录日志
		fmt.Printf("[%s] %s %s %d %v %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			clientIP,
			c.Request.Method,
			statusCode,
			duration,
			c.Request.URL.Path,
			userID,
		)

		// 记录到操作日志表（如果需要）
		// 这里可以调用操作日志服务记录详细操作
	}
}

// SecurityHeadersMiddleware 安全头部中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全相关的HTTP头
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件（处理panic）
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误
				fmt.Printf("Panic recovered: %v\n", err)

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})

				// 中止请求处理
				c.Abort()
			}
		}()

		c.Next()
	}
}
