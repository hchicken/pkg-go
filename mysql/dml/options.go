package dml

import (
	"github.com/hchicken/pkg-go/mysql"
)

// Option TODO
type Option func(*Options)

// Options TODO
// db查询options
type Options struct {
	Pool       mysql.Pool  // pool
	DbModel    interface{} // DB结构体
	ScanModel  interface{} // 查询结果
	Conditions interface{} // 查询条件
	In         []string    // in查询
	Like       []string    // like的查询条件
	Page       int         // 页码
	Limit      int         // 查询数量
	Total      *int64      // 查询数量
	SortField  string      // 排序
	STime      string      // 开始时间
	ETime      string      // 结束时间
	UpdateName string      // 更新key
	Values     []string    // 更新字段
}

func newOptions(opts ...Option) Options {
	opt := *new(Options)
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Pool TODO
func Pool(pool mysql.Pool) Option {
	return func(o *Options) {
		o.Pool = pool
	}
}

// DbModel TODO
// DB模型
func DbModel(model interface{}) Option {
	return func(o *Options) {
		o.DbModel = model
	}
}

// ScanModel 查询结果
func ScanModel(model interface{}) Option {
	return func(o *Options) {
		o.ScanModel = model
	}
}

// Conditions 查询条件
func Conditions(c interface{}) Option {
	return func(o *Options) {
		o.Conditions = c
	}
}

// In 查询
func In(in []string) Option {
	return func(o *Options) {
		o.In = in
	}
}

// Like TODO
// like查询
func Like(like []string) Option {
	return func(o *Options) {
		o.Like = like
	}
}

// Total 数据量
func Total(total *int64) Option {
	return func(o *Options) {
		o.Total = total
	}
}

// Limit 数据量
func Limit(limit int) Option {
	return func(o *Options) {
		o.Limit = limit
	}
}

// SortField 排序字段
func SortField(field string) Option {
	return func(o *Options) {
		o.SortField = field
	}
}

// STime 开始时间
func STime(t string) Option {
	return func(o *Options) {
		o.STime = t
	}
}

// ETime 结束时间
func ETime(t string) Option {
	return func(o *Options) {
		o.ETime = t
	}
}

// UpdateName 更新字段key
func UpdateName(name string) Option {
	return func(o *Options) {
		o.UpdateName = name
	}
}

// Values 更新的value
func Values(v []string) Option {
	return func(o *Options) {
		o.Values = v
	}
}
