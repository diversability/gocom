package redis

import "time"

func (p *Redis) Time() (res time.Time, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.Time().Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.Time().Result()
}
