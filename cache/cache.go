package cache

import (
	"fmt"
	"log"
)

// CachePool 缓存池接口
type CachePool interface {
	GetConnection() Connection // 获取链接
	Close() error              // 关闭连接
}

// NewCachePool 创建一个Redis客户端
func NewCachePool(opts ...Option) (CachePool, error) {
	client, err := newClient(opts...)
	if err != nil {
		log.Printf("Failed to create Redis client with options: %v", opts)
		return nil, fmt.Errorf("failed to create Redis client： %v", err)
	}
	return client, nil
}
