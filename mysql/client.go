package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type client struct {
	opts Options
	DB   *gorm.DB
}

func newClient(opts ...Option) (Pool, error) {
	options := newOptions(opts...)
	conn := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		options.user, options.password, options.uri, options.port, options.name,
	)
	// 临时解决
	cf := new(gorm.Config)
	pool, err := gorm.Open(mysql.Open(conn), cf)
	if err != nil {
		return nil, err
	}
	// 设置参数
	db, err := pool.DB()
	if err != nil {
		return nil, err
	}
	// 设置数据库连接池参数
	db.SetMaxOpenConns(options.MaxOpenConn)
	db.SetMaxIdleConns(options.MaxIdleConn)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	return &client{opts: options, DB: pool}, nil
}

// GetConn TODO
func (c client) GetConn() *gorm.DB {
	return c.DB
}
