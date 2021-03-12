package redis

func (p *Redis) LPush(key string, values ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.LPush(key, values...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.LPush(key, values...).Result()
}

func (p *Redis) RPush(key string, values ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.RPush(key, values...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.RPush(key, values...).Result()
}

func (p *Redis) LPop(key string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.LPop(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.LPop(key).Result()
}

func (p *Redis) RPop(key string) (res string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.RPop(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.RPop(key).Result()
}

func (p *Redis) LRange(key string, start, stop int64) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.LRange(key, start, stop).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.LRange(key, start, stop).Result()
}
