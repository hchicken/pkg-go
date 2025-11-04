package gormx

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBPool db连接池
type DBPool struct {
	opts Options
	DB   *gorm.DB
}

// NewDBPool get db pool
func NewDBPool(opts ...Option) (*DBPool, error) {
	options := newOptions(opts...)

	// use your own DB link if you set it up yourself
	if options.db != nil {
		return &DBPool{opts: options, DB: options.db}, nil
	}

	conn := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		options.user,
		options.password,
		options.uri,
		options.port,
		options.name,
	)
	cf := new(gorm.Config)
	pool, err := gorm.Open(mysql.Open(conn), cf)
	if err != nil {
		return nil, err
	}

	db, err := pool.DB()
	if err != nil {
		return nil, err
	}
	// 设置数据库连接池参数
	db.SetMaxOpenConns(options.MaxOpenConn)
	db.SetMaxIdleConns(options.MaxIdleConn)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	return &DBPool{opts: options, DB: pool}, nil
}

// GetConn get conn
func (c *DBPool) GetConn() *gorm.DB {
	if c.DB == nil {
		log.Fatal("database client not init")
	}
	return c.DB
}
