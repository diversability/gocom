package cache

import "github.com/diversability/gocom/cache/redis"

var redisCache *redis.Redis

// 初始化Redis 建立连接
func InitRedis(addrs []string, passwd string, db int, poolSize, minIdleConns, maxRetries int, cluster bool) {
	// 生成redis连接
	redisCache = redis.New(addrs, passwd, db, poolSize, minIdleConns, maxRetries, cluster)
}

// Redis 获取redis句柄
func Redis() *redis.Redis {
	return redisCache
}

// Close 关闭
func Close() {
	redisCache.Close()
}
