package models

import (
	"gorm.io/gorm"
)

const TableNameBrand = "brands"

// Brand mapped from table <roles>
type Brand struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string         `gorm:"column:name;not null" json:"name" binding:"required"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   *LocalTime     `json:"created_at"`
	UpdatedAt   *LocalTime     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // 软删除标志
}

// TableName Role's table name
func (*Brand) TableName() string {
	return TableNameBrand
}
