package model

// TabBaseModel 模型定义基类
type TabBaseModel struct {
	ID        uint     `json:"id" gorm:"primary_key"`
	CreatedBy string   `json:"created_by" gorm:"type:varchar(128);column:created_by;comment:'添加人'"`
	UpdatedBy string   `json:"updated_by" gorm:"type:varchar(128);column:updated_by;comment:'更新人'"`
	CreatedAt JsonTime `json:"created_at" gorm:"comment:'添加时间'"`
	UpdatedAt JsonTime `json:"updated_at" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP on update current_timestamp;comment:'更新时间'"`
}
