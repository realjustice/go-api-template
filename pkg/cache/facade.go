package cache

import (
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
)

// CacheFacade 缓存门面
type CacheFacade struct {
	manager cache.CacheInterface[string]
}

// NewCacheFacade 创建缓存门面
func NewCacheFacade(manager cache.CacheInterface[string]) *CacheFacade {
	return &CacheFacade{
		manager: manager,
	}
}

// Get 获取缓存
func (f *CacheFacade) Get(ctx context.Context, key string) (string, error) {
	value, err := f.manager.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return value, nil
}

// Set 设置缓存
func (f *CacheFacade) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return f.manager.Set(ctx, key, value, store.WithExpiration(ttl))
}

// Delete 删除缓存
func (f *CacheFacade) Delete(ctx context.Context, key string) error {
	return f.manager.Delete(ctx, key)
}

// Has 检查缓存是否存在
func (f *CacheFacade) Has(ctx context.Context, key string) bool {
	_, err := f.manager.Get(ctx, key)
	return err == nil
}

// Remember 记忆模式（Laravel 风格）
// 如果缓存存在则返回缓存值，否则执行回调函数并将结果缓存
func (f *CacheFacade) Remember(ctx context.Context, key string, ttl time.Duration, callback func() (string, error)) (string, error) {
	// 先尝试获取缓存
	value, err := f.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// 缓存未命中，执行回调
	value, err = callback()
	if err != nil {
		return "", err
	}

	// 存入缓存
	_ = f.Set(ctx, key, value, ttl)

	return value, nil
}

// Clear 清空所有缓存
func (f *CacheFacade) Clear(ctx context.Context) error {
	return f.manager.Clear(ctx)
}
