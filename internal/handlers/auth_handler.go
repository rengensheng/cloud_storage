package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"cloud-storage/internal/middleware"
	"cloud-storage/internal/models"
	"cloud-storage/internal/repositories"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo       *repositories.UserRepository
	authMiddleware *middleware.AuthMiddleware
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(
	userRepo *repositories.UserRepository,
	authMiddleware *middleware.AuthMiddleware,
) *AuthHandler {
	return &AuthHandler{
		userRepo:       userRepo,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册认证路由
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		auth.GET("/profile", h.RequireAuth(), h.GetProfile)
		auth.PUT("/profile", h.RequireAuth(), h.UpdateProfile)
		auth.PUT("/password", h.RequireAuth(), h.ChangePassword)
	}
}

// RequireAuth 要求认证的中间件包装
func (h *AuthHandler) RequireAuth() gin.HandlerFunc {
	return h.authMiddleware.Authenticate()
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	exists, err := (*h.userRepo).ExistsByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// 检查邮箱是否已存在
	exists, err = (*h.userRepo).ExistsByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check email"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// 哈希密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 设置用户角色
	role := models.RoleUser
	if req.Role != "" {
		role = req.Role
	}

	// 创建用户
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         role,
		IsActive:     true,
	}

	if err := (*h.userRepo).Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// 生成令牌
	accessToken, err := h.authMiddleware.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.authMiddleware.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    user.ToResponse(),
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600, // 1小时
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	user, err := (*h.userRepo).FindByUsername(req.Username)
	if err != nil {
		// 用户不存在或查询错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 检查用户是否活跃
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 更新最后登录时间
	if err := (*h.userRepo).UpdateLastLogin(user.ID); err != nil {
		// 记录错误但不影响登录
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// 生成令牌
	accessToken, err := h.authMiddleware.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.authMiddleware.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    user.ToResponse(),
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600, // 1小时
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取访问令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header is required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header format"})
		return
	}

	tokenString := parts[1]

	// 解析令牌获取过期时间
	claims, err := h.authMiddleware.ParseToken(tokenString)
	if err != nil {
		// 令牌无效，仍然返回成功
		c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
		return
	}

	// 将令牌加入黑名单
	expireTime := claims.ExpiresAt.Time
	if err := h.authMiddleware.BlacklistToken(tokenString, expireTime); err != nil {
		// 黑名单操作失败，记录错误但仍返回成功
		fmt.Printf("Failed to blacklist token: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 刷新令牌
	newAccessToken, newRefreshToken, err := h.authMiddleware.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed successfully",
		"tokens": gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600, // 1小时
		},
	})
}

// GetProfile 获取用户资料
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := (*h.userRepo).FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Username != nil {
		// 检查新用户名是否已存在
		exists, err := (*h.userRepo).ExistsByUsername(*req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
			return
		}
		if exists {
			// 检查是否是自己当前的用户名
			currentUser, err := (*h.userRepo).FindByID(userID)
			if err != nil || currentUser.Username != *req.Username {
				c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
				return
			}
		}
		updates["username"] = *req.Username
	}

	if req.Email != nil {
		// 检查新邮箱是否已存在
		exists, err := (*h.userRepo).ExistsByEmail(*req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check email"})
			return
		}
		if exists {
			// 检查是否是自己当前的邮箱
			currentUser, err := (*h.userRepo).FindByID(userID)
			if err != nil || currentUser.Email != *req.Email {
				c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
				return
			}
		}
		updates["email"] = *req.Email
	}

	if req.Role != nil {
		// 只有管理员可以修改角色
		userRole := c.MustGet("role").(string)
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to change role"})
			return
		}
		updates["role"] = *req.Role
	}

	if req.StorageQuota != nil {
		// 只有管理员可以修改存储配额
		userRole := c.MustGet("role").(string)
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to change storage quota"})
			return
		}
		updates["storage_quota"] = *req.StorageQuota
	}

	if req.IsActive != nil {
		// 只有管理员可以修改活跃状态
		userRole := c.MustGet("role").(string)
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to change active status"})
			return
		}
		updates["is_active"] = *req.IsActive
	}

	// 应用更新
	if len(updates) > 0 {
		if err := (*h.userRepo).Update(userID, updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}
	}

	// 获取更新后的用户信息
	user, err := (*h.userRepo).FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get updated profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    user.ToResponse(),
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户
	user, err := (*h.userRepo).FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	// 哈希新密码
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
		return
	}

	// 更新密码
	if err := (*h.userRepo).Update(userID, map[string]interface{}{
		"password_hash": string(newPasswordHash),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// ResetPassword 重置密码（需要邮箱验证）
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	// 重置密码功能需要邮箱服务
	c.JSON(http.StatusNotImplemented, gin.H{"error": "password reset not implemented yet"})
}

// VerifyEmail 验证邮箱
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	// 邮箱验证功能需要邮箱服务
	c.JSON(http.StatusNotImplemented, gin.H{"error": "email verification not implemented yet"})
}

// DeleteAccount 删除账户
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户
	user, err := (*h.userRepo).FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password is incorrect"})
		return
	}

	// 软删除用户账户
	if err := (*h.userRepo).SoftDelete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	// 获取访问令牌并加入黑名单
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := h.authMiddleware.ParseToken(tokenString)
			if err == nil {
				h.authMiddleware.BlacklistToken(tokenString, claims.ExpiresAt.Time)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
