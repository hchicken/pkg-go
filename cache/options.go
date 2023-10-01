package cache

import "time"

// Option 是一个函数类型，它接受一个指向Options的指针
type Option func(*Options)

// Options 是一个结构体，包含了Redis的配置选项
type Options struct {
	uri         string        // Redis服务器的URI
	db          int           // 使用的数据库编号
	password    string        // 连接Redis服务器的密码
	MaxIdle     int           // 连接池中最大空闲连接数
	MaxActive   int           // 连接池中最大活跃连接数
	IdleTimeout time.Duration // 连接池中连接的最大空闲时间
}

// newOptions 创建一个新的Options实例，并应用提供的选项
func newOptions(opts ...Option) Options {
	opt := *new(Options)
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Uri 是一个Option，用于设置Redis服务器的URI
func Uri(uri string) Option {
	return func(o *Options) {
		o.uri = uri
	}
}

// DB 是一个Option，用于设置使用的数据库编号
func DB(db int) Option {
	return func(o *Options) {
		o.db = db
	}
}

// PassWord 是一个Option，用于设置连接Redis服务器的密码
func PassWord(pwd string) Option {
	return func(o *Options) {
		o.password = pwd
	}
}

// MaxIdle 是一个Option，用于设置连接池中最大空闲连接数
func MaxIdle(num int) Option {
	return func(o *Options) {
		o.MaxIdle = num
	}
}

// MaxActive 是一个Option，用于设置连接池中最大活跃连接数
func MaxActive(num int) Option {
	return func(o *Options) {
		o.MaxActive = num
	}
}

// IdleTimeout 是一个Option，用于设置连接池中连接的最大空闲时间
func IdleTimeout(num time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = num
	}
}
