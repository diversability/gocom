package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

func (p *Redis) ZRevRange(key string, start, stop int64) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZRevRange(key, start, stop).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZRevRange(key, start, stop).Result()
}

func (p *Redis) ZRevRangeByScore(key string, max, min string, offset, count int64) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZRevRangeByScore(key,
			redis.ZRangeBy{
				Min:    min,
				Max:    max,
				Offset: offset,
				Count:  count,
			}).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZRevRangeByScore(key,
		redis.ZRangeBy{
			Min:    min,
			Max:    max,
			Offset: offset,
			Count:  count,
		}).Result()
}

func (p *Redis) ZAdd(key string, score float64, member interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZAdd(key,
			redis.Z{
				Score:  score,
				Member: member,
			}).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZAdd(key,
		redis.Z{
			Score:  score,
			Member: member,
		}).Result()
}

func (p *Redis) ZCount(key, min, max string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZCount(key, min, max).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZCount(key, min, max).Result()
}

func (p *Redis) ZScore(key, member string) (res float64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZScore(key, member).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZScore(key, member).Result()
}

func (p *Redis) ZRange(key string, start, stop int64) (res []string, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZRange(key, start, stop).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZRange(key, start, stop).Result()
}

func (p *Redis) ZIncrBy(key string, incr float64, member string) (res float64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZIncrBy(key, incr, member).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZIncrBy(key, incr, member).Result()
}

func (p *Redis) ZRem(key string, members ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZRem(key, members...).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZRem(key, members...).Result()
}

func (p *Redis) ZUnionStore(dest string, weights []float64, aggregate string, keys ...string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		return 0, fmt.Errorf("cluster unsupport <zunionstore>.")
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZUnionStore(dest,
		redis.ZStore{
			Weights:   weights,
			Aggregate: aggregate,
		},
		keys...).Result()
}

func (p *Redis) ZCard(key string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZCard(key).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZCard(key).Result()
}

func (p *Redis) ZRemRangeByScore(key, min, max string) (res int64, err error) {
	// 安全检查
	if p == nil {
		return res, errRedis
	}

	if p.cluster {
		// 客户端安全检查
		if p.cclient == nil {
			return res, errRedisCClient
		}

		return p.cclient.ZRemRangeByScore(key, min, max).Result()
	}

	// 客户端安全检查
	if p.client == nil {
		return res, errRedisClient
	}

	return p.client.ZRemRangeByScore(key, min, max).Result()
}
