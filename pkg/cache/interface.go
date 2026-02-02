package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Has 检查缓存是否存在
	Has(ctx context.Context, key string) bool

	// Remember 记忆模式（缓存未命中时执行回调并缓存结果）
	Remember(ctx context.Context, key string, ttl time.Duration, callback func() (string, error)) (string, error)

	// Clear 清空所有缓存
	Clear(ctx context.Context) error
}
