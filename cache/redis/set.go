package redis

import "fmt"

func (p *Redis) SCard(key string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SCard(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SCard(key).Result()
}

func (p *Redis) SAdd(key string, members ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SAdd(key, members...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SAdd(key, members...).Result()
}

func (p *Redis) SMembers(key string) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SMembers(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SMembers(key).Result()
}

func (p *Redis) SIsmember(key string, member interface{}) (res bool, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SIsMember(key, member).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SIsMember(key, member).Result()
}

func (p *Redis) SInter(keys ...string) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return nil, fmt.Errorf("cluster unsupport <sinter>.")
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SInter(keys...).Result()
}

func (p *Redis) SPop(key string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SPop(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SPop(key).Result()
}

func (p *Redis) SRem(key string, members ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.SRem(key, members...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.SRem(key, members...).Result()
}
