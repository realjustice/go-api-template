package redis

import (
	"context"
	"fmt"
	"time"

	"go-api-template/pkg/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(cfg *config.Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return &Client{Client: client}, nil
}

// Close 关闭 Redis 连接
func (c *Client) Close() error {
	return c.Client.Close()
}
