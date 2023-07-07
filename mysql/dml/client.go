package dml

import (
	"fmt"
	"sync"

	"github.com/hchicken/pkg-go/mysql"
	"github.com/hchicken/pkg-go/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type client struct {
	pool mysql.Pool
	opts Options
	once sync.Once
}

func newClient(opts ...Option) DB {
	cli := new(client)
	options := newOptions(opts...)
	cli.opts = options
	cli.pool = options.Pool
	return cli
}

// BaseRead 基础查询
func (c *client) BaseRead() (*gorm.DB, error) {
	var conditions map[string]interface{}
	err := util.StructDecode(c.opts.Conditions, &conditions)
	if err != nil {
		return nil, err
	}
	// 去掉为空的值
	for key, value := range conditions {
		if value == "" {
			delete(conditions, key)
		}
	}
	pool := c.pool.GetConn()
	// like查询
	for _, v := range c.opts.Like {
		value := conditions[v]
		if value == nil {
			continue
		}
		sql := fmt.Sprintf("%%%v%%", conditions[v])
		pool = pool.Where(fmt.Sprintf("`%v` LIKE  ?", v), sql)
		delete(conditions, v)
	}

	// 遍历查询
	for _, v := range c.opts.In {
		value := conditions[v]
		if value == nil {
			continue
		}
		sql := fmt.Sprintf("%v in  (?)", v)
		pool = pool.Where(sql, value)
		delete(conditions, v)
	}

	// 直接查询
	pool = pool.Where(conditions)

	// 时间查询
	if c.opts.STime != "" && c.opts.ETime != "" {
		pool = pool.Where("created_at BETWEEN ? AND ?", c.opts.STime, c.opts.ETime)
	}
	return pool, err
}

// Read 查询
func (c *client) Read() error {
	pool, err := c.BaseRead()
	if err != nil {
		return err
	}
	// 分页
	if c.opts.Limit != 0 {
		pool = pool.Limit(c.opts.Limit)
	}
	if c.opts.Limit != 0 && c.opts.Page != 0 {
		pool = pool.Offset((c.opts.Page - 1) * c.opts.Limit)
	}

	if c.opts.SortField == "" {
		c.opts.SortField = "id DESC"
	}
	return pool.Model(c.opts.DbModel).Order(c.opts.SortField).Scan(c.opts.ScanModel).Error
}

// Create 添加
func (c *client) Create() error {
	err := c.pool.GetConn().Create(c.opts.DbModel).Error
	return err
}

// CreateOrUpdate 添加或者更新
func (c *client) CreateOrUpdate() error {
	conn := c.pool.GetConn()
	err := conn.Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: c.opts.UpdateName,
		}},
		DoUpdates: clause.AssignmentColumns(c.opts.Values),
	}).Create(c.opts.DbModel).Error
	return err
}

// Update 更新
func (c *client) Update() error {
	return nil
}

// Delete 删除
func (c *client) Delete() error {
	pool, err := c.BaseRead()
	if err != nil {
		return err
	}
	return pool.Delete(c.opts.DbModel).Error
}
