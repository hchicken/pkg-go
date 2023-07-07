package proto

// ReqQueryBase 查询请求
type ReqQueryBase struct {
	Limit *int   `json:"limit,omitempty" form:"limit,omitempty"`                        // 数量
	Page  *int   `json:"page,omitempty" form:"page,omitempty"`                          // 分页
	Sort  string `json:"sort,omitempty" form:"sort,omitempty"`                          // 排序
	STime string `json:"s_time,omitempty" form:"s_time,omitempty" binding:"timeString"` // 数据添加时间
	ETime string `json:"e_time,omitempty" form:"e_time,omitempty" binding:"timeString"` // 添加数据时间
}
