package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// Conn TODO
type Conn redis.Conn

// Client TODO
type Client struct {
	opts Options
	pool *redis.Pool
}

func newClient(opts ...Option) CachePool {
	cli := new(Client)
	options := newOptions(opts...)
	cli.opts = options
	cli.pool = &redis.Pool{
		MaxIdle:     options.MaxIdle,     // 最大空闲数
		IdleTimeout: options.IdleTimeout, // 最大空闲时间
		MaxActive:   options.MaxActive,   // 最大数
		Dial: func() (redis.Conn, error) {
			dbOption := redis.DialDatabase(options.db)
			pwOption := redis.DialPassword(options.password)

			c, err := redis.Dial("tcp", options.uri, dbOption, pwOption)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return cli
}

// GetCli 获取redis客户端
func (cli *Client) GetCli() Conn {
	return cli.pool.Get()
}
