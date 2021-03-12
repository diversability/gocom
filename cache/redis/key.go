package redis

import (
	"time"
)

func (p *Redis) Keys(key string) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Keys(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Keys(key).Result()
}

func (p *Redis) Expire(key string, seconds int64) (res bool, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Expire(key, time.Duration(seconds)*time.Second).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Expire(key, time.Duration(seconds)*time.Second).Result()
}

func (p *Redis) Exists(keys ...string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Exists(keys...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Exists(keys...).Result()
}

func (p *Redis) Del(keys ...string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Del(keys...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Del(keys...).Result()
}

func (p *Redis) IncrBy(key string, value int64) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.IncrBy(key, value).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.IncrBy(key, value).Result()
}

func (p *Redis) Rename(key string, newkey string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Rename(key, newkey).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Rename(key, newkey).Result()
}
