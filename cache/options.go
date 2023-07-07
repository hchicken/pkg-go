package cache

import "time"

// Option TODO
type Option func(*Options)

// Options TODO
// redis options
type Options struct {
	uri         string
	db          int
	password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

func newOptions(opts ...Option) Options {
	opt := *new(Options)
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Uri 设置host
func Uri(uri string) Option {
	return func(o *Options) {
		o.uri = uri
	}
}

// DB 设置DB
func DB(db int) Option {
	return func(o *Options) {
		o.db = db
	}
}

// PassWord 设置密码
func PassWord(pwd string) Option {
	return func(o *Options) {
		o.password = pwd
	}
}

// MaxIdle 最大空闲数
func MaxIdle(num int) Option {
	return func(o *Options) {
		o.MaxIdle = num
	}
}

// MaxActive 最大数连接数
func MaxActive(num int) Option {
	return func(o *Options) {
		o.MaxActive = num
	}
}

// IdleTimeout 最大空闲时间
func IdleTimeout(num time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = num
	}
}
