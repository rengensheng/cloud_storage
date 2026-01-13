package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud-storage/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		PoolSize:    100,
		PoolTimeout: 30 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = redisClient
	log.Println("Redis connection established successfully")
	return redisClient, nil
}

// GetRedis 获取Redis客户端实例
func GetRedis() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// Redis操作封装

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Del 删除键
func Del(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// HSet 设置哈希字段
func HSet(key string, values ...interface{}) error {
	return RedisClient.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func HGet(key, field string) (string, error) {
	return RedisClient.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func HGetAll(key string) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, key).Result()
}

// LPush 列表左推入
func LPush(key string, values ...interface{}) error {
	return RedisClient.LPush(ctx, key, values...).Err()
}

// RPop 列表右弹出
func RPop(key string) (string, error) {
	return RedisClient.RPop(ctx, key).Result()
}

// SAdd 集合添加成员
func SAdd(key string, members ...interface{}) error {
	return RedisClient.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	return RedisClient.SMembers(ctx, key).Result()
}

// ZAdd 有序集合添加成员
func ZAdd(key string, members ...redis.Z) error {
	return RedisClient.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func ZRange(key string, start, stop int64) ([]string, error) {
	return RedisClient.ZRange(ctx, key, start, stop).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return RedisClient.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	return RedisClient.Decr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(key string) (time.Duration, error) {
	return RedisClient.TTL(ctx, key).Result()
}

// Redis使用场景示例

// CacheUserToken 缓存用户令牌
func CacheUserToken(userID string, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user:token:%s", userID)
	return Set(key, token, expiration)
}

// GetUserToken 获取用户令牌
func GetUserToken(userID string) (string, error) {
	key := fmt.Sprintf("user:token:%s", userID)
	return Get(key)
}

// RemoveUserToken 移除用户令牌
func RemoveUserToken(userID string) error {
	key := fmt.Sprintf("user:token:%s", userID)
	return Del(key)
}

// RateLimit 速率限制
func RateLimit(key string, limit int, window time.Duration) (bool, error) {
	current, err := Incr(key)
	if err != nil {
		return false, err
	}

	if current == 1 {
		Expire(key, window)
	}

	return current <= int64(limit), nil
}

// CacheFileMetadata 缓存文件元数据
func CacheFileMetadata(fileID string, metadata map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("file:metadata:%s", fileID)

	// 将map转换为interface{}切片
	var values []interface{}
	for k, v := range metadata {
		values = append(values, k, v)
	}

	return HSet(key, values...)
}

// GetFileMetadata 获取文件元数据
func GetFileMetadata(fileID string) (map[string]string, error) {
	key := fmt.Sprintf("file:metadata:%s", fileID)
	return HGetAll(key)
}
