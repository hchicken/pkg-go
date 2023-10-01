package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisConnection 是 redis.Conn 的别名，用于表示一个 Redis 连接
type Connection redis.Conn

// RedisClient 是一个包含连接池和配置选项的 Redis 客户端
type RedisClient struct {
	options Options     // Redis 连接的配置选项
	pool    *redis.Pool // Redis 连接池
}

// createPool 创建一个 Redis 连接池
func createPool(options Options) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     options.MaxIdle,     // 连接池中的最大空闲连接数
		IdleTimeout: options.IdleTimeout, // 空闲连接的最大等待时间
		MaxActive:   options.MaxActive,   // 连接池中的最大活动连接数
		Dial: func() (redis.Conn, error) { // 创建新的 Redis 连接的函数
			dbOption := redis.DialDatabase(options.db)
			pwOption := redis.DialPassword(options.password)

			connection, err := redis.Dial("tcp", options.uri, dbOption, pwOption)
			if err != nil {
				// 如果连接失败，记录错误信息
				return nil, err
			}
			return connection, err
		},
		TestOnBorrow: func(connection redis.Conn, t time.Time) error { // 从连接池借用连接时的测试函数
			_, err := connection.Do("PING")
			return err
		},
	}
}

// newClient 创建一个新的 Redis 客户端
func newClient(options ...Option) CachePool {
	client := new(RedisClient)
	clientOptions := newOptions(options...)
	client.options = clientOptions
	client.pool = createPool(clientOptions)
	return client
}

// GetConnection 从连接池中获取一个 Redis 连接
func (client *RedisClient) GetConnection() Connection {
	return client.pool.Get()
}
