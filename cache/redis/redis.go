package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	cluster bool
	client  *redis.Client
	cclient *redis.ClusterClient
}

var (
	errRedis        = fmt.Errorf("redis is nil")
	errRedisClient  = fmt.Errorf("redis client is nil")
	errRedisCClient = fmt.Errorf("redis cclient is nil")
)

func New(addrs []string, passwd string, db int,
	poolSize, minIdleConns, maxRetries int, cluster bool) *Redis {
	if len(addrs) == 0 {
		panic("Redis's addrs is empty.")
	}

	r := &Redis{cluster: cluster}

	if cluster {
		r.cclient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addrs,
			Password:     passwd,
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
			MaxRetries:   maxRetries,
		})

		// check if redis server is ok.
		if _, err := r.cclient.Ping().Result(); err != nil {
			panic(err)
		}
	} else {
		r.client = redis.NewClient(&redis.Options{
			Addr:         addrs[0],
			Password:     passwd,
			DB:           db,
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
			MaxRetries:   maxRetries,
		})

		// check if redis server is ok.
		if _, err := r.client.Ping().Result(); err != nil {
			panic(err)
		}
	}

	return r
}

// Close 关闭redis
func (p *Redis) Close() {
	if p != nil {
		if p.cluster {
			if p.cclient != nil {
				p.cclient.Close()
			}
		} else {
			if p.client != nil {
				p.client.Close()
			}
		}
	}
}
