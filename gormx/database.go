package gormx

import (
	"fmt"
	"reflect"

	"github.com/hchicken/pkg-go/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Database represents a database connection and operations
type Database struct {
	db   *gorm.DB
	opts ConnectionOptions
}

// NewDatabase creates a new Database instance
func NewDatabase(opts ...ConnectionOption) *Database {
	options := newConnectionOptions(opts...)
	return &Database{
		db:   options.Pool,
		opts: options,
	}
}

func (d *Database) decodeAndCleanConditions() (map[string]interface{}, error) {
	var conditions map[string]interface{}
	err := util.StructDecode(d.opts.Conditions, &conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to decode conditions: %w", err)
	}

	// Remove specified fields
	for _, field := range d.opts.ExcludeFields {
		delete(conditions, field)
	}

	return conditions, nil
}

// prepareQuery sets up the base query with conditions
func (d *Database) prepareQuery() (*gorm.DB, error) {
	conditions, err := d.decodeAndCleanConditions()
	if err != nil {
		return nil, fmt.Errorf("failed to decode conditions: %w", err)
	}

	query := d.db.Model(d.opts.DbModel)

	// Apply LIKE conditions
	for _, field := range d.opts.Like {
		if value, ok := conditions[field]; ok && value != "" {
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", value))
			delete(conditions, field)
		}
	}

	// Apply IN conditions
	for _, field := range d.opts.In {
		if value, ok := conditions[field]; ok && value != nil {
			query = query.Where(fmt.Sprintf("%s IN (?)", field), value)
			delete(conditions, field)
		}
	}

	// Apply remaining conditions
	query = query.Where(conditions)

	// Apply time range if specified
	if d.opts.StartTime != "" && d.opts.EndTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", d.opts.StartTime, d.opts.EndTime)
	}

	return query, nil
}

// applyPagination adds limit and offset to the query
func (d *Database) applyPagination(query *gorm.DB) *gorm.DB {
	if d.opts.Limit != 0 {
		query = query.Limit(d.opts.Limit)
	}
	if d.opts.Offset != 0 {
		query = query.Offset(d.opts.Offset)
	} else if d.opts.Page != 0 {
		query = query.Offset((d.opts.Page - 1) * d.opts.Limit)
	}
	return query
}

// applyOrder adds Order to the query
func (d *Database) applyOrder(query *gorm.DB) *gorm.DB {
	if d.opts.SortField == "" {
		d.opts.SortField = "id DESC"
	}
	query = query.Order(d.opts.SortField)
	return query
}

// executeQuery prepares and executes a query with the given operation
func (d *Database) executeQuery(operation func(*gorm.DB) *gorm.DB) error {
	query, err := d.prepareQuery()
	if err != nil {
		return err
	}

	if d.opts.Total != nil {
		if err := query.Count(d.opts.Total).Error; err != nil {
			return fmt.Errorf("failed to count total: %w", err)
		}
	}

	query = d.applyPagination(query)
	query = d.applyOrder(query)

	if d.opts.Debug {
		query = query.Debug()
	}

	return operation(query).Error
}

// Query retrieves records based on the set options
func (d *Database) Query() error {
	return d.executeQuery(func(query *gorm.DB) *gorm.DB {
		return query.Scan(d.opts.ScanModel)
	})
}

// First retrieves the first record that matches the query
func (d *Database) First() error {
	return d.executeQuery(func(query *gorm.DB) *gorm.DB {
		return query.First(d.opts.ScanModel)
	})
}

// Create adds a new record to the database
func (d *Database) Create(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer to a struct")
	}
	return d.db.Create(value).Error
}

// CreateOrUpdate adds a new record or updates an existing one
func (d *Database) CreateOrUpdate() error {
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: d.opts.UpdateName}},
		DoUpdates: clause.AssignmentColumns(d.opts.Values),
	}).Create(d.opts.DbModel).Error
}

// Update modifies existing records (placeholder for future implementation)
func (d *Database) Update() error {
	// TODO: Implement update logic
	return nil
}

// Delete removes records based on the set conditions
func (d *Database) Delete() error {
	query, err := d.prepareQuery()
	if err != nil {
		return err
	}
	return query.Delete(d.opts.DbModel).Error
}
