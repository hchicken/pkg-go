package gormx

import "gorm.io/gorm"

// ConnectionOption ...
type ConnectionOption func(*ConnectionOptions)

// db查询ConnectionOptions
type ConnectionOptions struct {
	Pool          *gorm.DB    // pool
	DbModel       interface{} // DB结构体
	ScanModel     interface{} // 查询结果
	Conditions    interface{} // 查询条件
	ExcludeFields []string    // 不查询的字段
	In            []string    // in查询
	Like          []string    // like的查询条件
	Page          int         // 页码
	Limit         int         // 查询数量
	Offset        int         // 偏移量
	Total         *int64      // 总数
	SortField     string      // 排序
	StartTime     string      // 开始时间
	EndTime       string      // 结束时间
	UpdateName    string      // 更新key
	Values        []string    // 更新字段

	Debug bool // 是否debug查询
}

func newConnectionOptions(opts ...ConnectionOption) ConnectionOptions {
	opt := ConnectionOptions{
		ExcludeFields: []string{"limit", "page", "sort", "s_time", "e_time"},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Pool ...
func WithConnPool(pool *gorm.DB) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Pool = pool
	}
}

// DbModel DB模型
func WithConnDbModel(model interface{}) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.DbModel = model
	}
}

// ScanModel 查询结果
func WithConnScanModel(model interface{}) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.ScanModel = model
	}
}

// Conditions 查询条件
func WithConnConditions(c interface{}) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Conditions = c
	}
}

// WithConnExcludeFields ...
func WithConnExcludeFields(fields []string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.ExcludeFields = fields
	}
}

// In 查询
func WithConnIn(in []string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.In = in
	}
}

// Like like查询
func WithConnLike(like []string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Like = like
	}
}

// Total 数据量
func WithConnTotal(total *int64) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Total = total
	}
}

// Limit 数据量
func WithConnLimit(limit int) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Limit = limit
	}
}

// Offset 偏移量
func WithConnOffset(offset int) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Offset = offset
	}
}

// Page 页码
func WithConnPage(page int) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Page = page
	}
}

// SortField 排序字段
func WithConnSortField(field string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.SortField = field
	}
}

// STime 开始时间
func WithConnStartTime(t string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.StartTime = t
	}
}

// ETime 结束时间
func WithConnEndTime(t string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.EndTime = t
	}
}

// UpdateName 更新字段key
func WithConnUpdateName(name string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.UpdateName = name
	}
}

// Values 更新的value
func WithConnValues(v []string) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Values = v
	}
}

// Values 更新的value
func WithConnDebug(b bool) ConnectionOption {
	return func(o *ConnectionOptions) {
		o.Debug = b
	}
}
