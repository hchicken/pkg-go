package cache

// CachePool 缓存池
type CachePool interface {
	GetConnection() Connection // 获取链接
}

// NewCachePool 创建一个redis客户端
func NewCachePool(opts ...Option) CachePool {
	return newClient(opts...)
}
