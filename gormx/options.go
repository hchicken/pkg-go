package gormx

import "time"

// Option ...
type Option func(*Options)

// Options gormx options
type Options struct {
	uri             string
	port            string
	name            string
	user            string
	password        string
	MaxOpenConn     int           // 最大连接
	MaxIdleConn     int           // 最大空闲连接
	ConnMaxLifetime time.Duration // 最大空闲时间
}

func newOptions(opts ...Option) Options {
	opt := *new(Options)
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Uri ...
func Uri(uri string) Option {
	return func(o *Options) {
		o.uri = uri
	}
}

// Name ...
func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

// Port ...
func Port(port string) Option {
	return func(o *Options) {
		o.port = port
	}
}

// User ...
func User(user string) Option {
	return func(o *Options) {
		o.user = user
	}
}

// PassWord ...
func PassWord(pwd string) Option {
	return func(o *Options) {
		o.password = pwd
	}
}

// MaxOpenConn 最大连接数
func MaxOpenConn(n int) Option {
	return func(o *Options) {
		o.MaxOpenConn = n
	}
}

// MaxIdleConn 最大空闲数
func MaxIdleConn(n int) Option {
	return func(o *Options) {
		o.MaxIdleConn = n
	}
}

// ConnMaxLifetime 最大存活时间
func ConnMaxLifetime(t time.Duration) Option {
	return func(o *Options) {
		o.ConnMaxLifetime = t
	}
}
