package redis

import (
	"time"
)

func (p *Redis) Get(key string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Get(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Get(key).Result()
}

func (p *Redis) Set(key string, value interface{}) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Set(key, value, 0).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Set(key, value, 0).Result()
}

func (p *Redis) SetEx(key string, seconds int64, value interface{}) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Set(key, value, time.Duration(seconds)*time.Second).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Set(key, value, time.Duration(seconds)*time.Second).Result()
}

func (p *Redis) MGet(keys ...string) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	vals := make([]interface{}, len(keys))
	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		if vals, err = p.cclient.MGet(keys...).Result(); err != nil {
			return res, err
		}

		// 转为[]string
		return strings(vals, err)
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	if vals, err = p.client.MGet(keys...).Result(); err != nil {
		return res, err
	}

	// 转为[]string
	return strings(vals, err)
}

func (p *Redis) SetBit(key string, offset int64, value int) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SetBit(key, offset, value).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SetBit(key, offset, value).Result()
}

func (p *Redis) GetBit(key string, offset int64) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.GetBit(key, offset).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.GetBit(key, offset).Result()
}

func (p *Redis) SetNx(key string, value interface{}) (res bool, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SetNX(key, value, 0).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SetNX(key, value, 0).Result()
}

func (p *Redis) SetNxEx(key string, value interface{}, seconds int64) (res bool, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SetNX(key, value, time.Duration(seconds)*time.Second).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SetNX(key, value, time.Duration(seconds)*time.Second).Result()
}

func (p *Redis) Append(key, value string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Append(key, value).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Append(key, value).Result()
}
