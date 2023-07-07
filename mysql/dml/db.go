package dml

// DB 接口
type DB interface {
	Read() error
	Create() error
	CreateOrUpdate() error
	Update() error
	Delete() error
}

// NewDB 获取db客户端
func NewDB(opts ...Option) DB {
	return newClient(opts...)
}
