package model

import "time"

// Demo 演示模型
type Demo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"type:varchar(200);not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Status    int       `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Demo) TableName() string {
	return "demos"
}
