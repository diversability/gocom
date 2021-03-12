package redis

import (
	"github.com/go-redis/redis"
)

func (p *Redis) HMGet(key string, fields ...string) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	vals := make([]interface{}, len(fields))
	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		if vals, err = p.cclient.HMGet(key, fields...).Result(); err != nil {
			return res, err
		}

		// 转为[]string
		return strings(vals, err)
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	if vals, err = p.client.HMGet(key, fields...).Result(); err != nil {
		return res, err
	}

	// 转为[]string
	return strings(vals, err)
}

func (p *Redis) HGet(key, field string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.HGet(key, field).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.HGet(key, field).Result()
}

func (p *Redis) HGetAll(key string) (res map[string]string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.HGetAll(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.HGetAll(key).Result()
}

func (p *Redis) HMSet(key string, fields map[string]interface{}) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.HMSet(key, fields).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.HMSet(key, fields).Result()
}

func (p *Redis) HSet(key, field string, value interface{}) (res bool, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.HSet(key, field, value).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.HSet(key, field, value).Result()
}

func (p *Redis) HIncrBy(key, field string, incr int64) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.HIncrBy(key, field, incr).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.HIncrBy(key, field, incr).Result()
}

func (p *Redis) BatchHGetAll(keys ...string) (res map[string]map[string]string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	cmds := make(map[string]*redis.StringStringMapCmd, len(keys))
	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		if _, err := p.cclient.Pipelined(func(pipe redis.Pipeliner) error {
			for _, key := range keys {
				cmds[key] = pipe.HGetAll(key)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		// 客户端安全检查
		if p.client == nil {
			return res, errRedisClient
		}

		if _, err := p.client.Pipelined(func(pipe redis.Pipeliner) error {
			for _, key := range keys {
				cmds[key] = pipe.HGetAll(key)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	result := make(map[string]map[string]string, len(cmds))
	for _, key := range keys {
		if v, ok := cmds[key]; ok && v != nil {
			result[key] = v.Val()
		}
	}

	return result, nil
}
