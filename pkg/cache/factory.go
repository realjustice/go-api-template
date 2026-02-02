package cache

import (
	"fmt"
	"time"

	"go-api-template/pkg/config"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

// CacheDriver 缓存驱动类型
type CacheDriver string

const (
	DriverRedis  CacheDriver = "redis"
	DriverMemory CacheDriver = "memory"
)

// NewCacheManager 根据配置创建缓存管理器
func NewCacheManager(cfg *config.Config, redisClient *redis.Client) (cache.CacheInterface[string], error) {
	driver := CacheDriver(cfg.Cache.Driver)

	switch driver {
	case DriverRedis:
		if redisClient == nil {
			return nil, fmt.Errorf("redis client is required for redis driver")
		}
		redisStore := redis_store.NewRedis(redisClient)
		return cache.New[string](redisStore), nil

	case DriverMemory:
		// 使用配置的 TTL 作为默认过期时间
		defaultTTL := time.Duration(cfg.Cache.TTL) * time.Second
		gocacheClient := gocache.New(defaultTTL, defaultTTL*2)
		gocacheStore := gocache_store.NewGoCache(gocacheClient)
		return cache.New[string](gocacheStore), nil

	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", driver)
	}
}

// NewChainCache 创建多级缓存（L1: Memory, L2: Redis）
// 先查内存缓存（快），未命中再查 Redis
func NewChainCache(cfg *config.Config, redisClient *redis.Client) (cache.CacheInterface[string], error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is required for chain cache")
	}

	// L1: 内存缓存（快）
	defaultTTL := time.Duration(cfg.Cache.TTL) * time.Second
	memoryStore := gocache_store.NewGoCache(
		gocache.New(defaultTTL, defaultTTL*2),
	)

	// L2: Redis 缓存（持久）
	redisStore := redis_store.NewRedis(redisClient)

	// 创建链式缓存
	chainCache := cache.NewChain[string](
		cache.New[string](memoryStore),
		cache.New[string](redisStore),
	)

	return chainCache, nil
}
