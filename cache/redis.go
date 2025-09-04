package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Connection 是 redis.Conn 的别名，用于表示一个 Redis 连接
type Connection redis.Conn

// RedisClient 是一个包含连接池和配置选项的 Redis 客户端
type RedisClient struct {
	options Options     // Redis 连接的配置选项
	pool    *redis.Pool // Redis 连接池
}

// 编译时检查RedisClient是否实现CachePool接口
var _ CachePool = (*RedisClient)(nil)

// createPool 创建一个 Redis 连接池
func createPool(options Options) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     options.MaxIdle,     // 连接池中的最大空闲连接数
		IdleTimeout: options.IdleTimeout, // 空闲连接的最大等待时间
		MaxActive:   options.MaxActive,   // 连接池中的最大活动连接数
		Dial: func() (redis.Conn, error) { // 创建新的 Redis 连接的函数
			// 设置连接选项
			dialOptions := []redis.DialOption{
				redis.DialDatabase(options.db),
				redis.DialConnectTimeout(5 * time.Second), // 连接超时
				redis.DialReadTimeout(3 * time.Second),    // 读超时
				redis.DialWriteTimeout(3 * time.Second),   // 写超时
			}

			// 如果有密码，添加密码选项
			if options.password != "" {
				dialOptions = append(dialOptions, redis.DialPassword(options.password))
			}

			connection, err := redis.Dial("tcp", options.uri, dialOptions...)
			if err != nil {
				log.Printf("Failed to connect to Redis: %v, URI: %s, DB: %d", err, options.uri, options.db)
				return nil, err
			}
			return connection, nil
		},
		TestOnBorrow: func(connection redis.Conn, t time.Time) error { // 从连接池借用连接时的测试函数
			if connection == nil {
				return fmt.Errorf("connection is nil")
			}
			_, err := connection.Do("PING")
			if err != nil {
				log.Printf("Redis connection test failed: %v", err)
			}
			return err
		},
	}

	// 测试连接池是否能正常工作
	testConn := pool.Get()
	defer func() {
		if testConn != nil {
			if err := testConn.Close(); err != nil {
				panic(fmt.Sprintf("Failed to close Redis test connection: %v", err))
			}
		}
	}()

	if testConn.Err() != nil {
		return nil, fmt.Errorf("failed to get test connection from pool: %v", testConn.Err())
	}

	// 执行一次PING测试
	_, err := testConn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("Redis connection test failed: %v", err)
	}

	log.Printf("Redis connection pool created successfully, URI: %s, DB: %d", options.uri, options.db)
	return pool, nil
}

// newClient 创建一个新的 Redis 客户端
func newClient(options ...Option) (CachePool, error) {
	client := new(RedisClient)
	clientOptions := newOptions(options...)
	client.options = clientOptions
	pool, err := createPool(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis pool: %v", err)
	}
	client.pool = pool
	return client, nil
}

// GetConnection 从连接池中获取一个 Redis 连接
func (client *RedisClient) GetConnection() Connection {
	return client.pool.Get()
}

// Close 释放连接池资源
func (client *RedisClient) Close() error {
	return client.pool.Close()
}
