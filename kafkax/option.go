package kafkax

import (
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Option func(*Options)

// Options ...
type Options struct {
	address   []string
	username  string
	password  string
	sasl      sasl.Mechanism
	algorithm scram.Algorithm
}

// newOptions ...
func newOptions(opts ...Option) Options {
	options := Options{
		algorithm: scram.SHA512, // 支持sha256和sha512
	}
	for _, o := range opts {
		o(&options)
	}

	// 账号密码不为空的时候设置
	if options.username != "" && options.password != "" {
		// 设置 algo and sasl
		switch options.algorithm {
		case scram.SHA512:
			mechanism, err := scram.Mechanism(options.algorithm, options.username, options.password)
			if err != nil {
				panic(err)
			}
			options.sasl = mechanism
		default:
			mechanism := plain.Mechanism{
				Username: options.username,
				Password: options.password,
			}
			options.sasl = mechanism
		}
	}

	return options
}

// Address ...
func Address(address []string) Option {
	return func(o *Options) {
		o.address = address
	}
}

// User ...
func User(username string) Option {
	return func(o *Options) {
		o.username = username
	}
}

// Password ...
func Password(password string) Option {
	return func(o *Options) {
		o.password = password
	}
}

// Algorithm ...
func Algorithm(algorithm scram.Algorithm) Option {
	return func(o *Options) {
		o.algorithm = algorithm
	}
}
