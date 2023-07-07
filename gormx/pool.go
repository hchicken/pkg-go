package gormx

import (
	"gorm.io/gorm"
)

// Pool db pool
type Pool interface {
	GetConn() *gorm.DB
}

// NewDBPool 获取一个db pool
func NewDBPool(opts ...Option) (Pool, error) {
	return newClient(opts...)
}
