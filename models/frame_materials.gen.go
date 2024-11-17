// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"gorm.io/gorm"
)

const TableNameFrameMaterial = "frame_materials"

// FrameMaterial mapped from table <frame_materials>
type FrameMaterial struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name;not null" json:"name" binding:"required"`
	Description string `gorm:"column:description;" json:"description"`

	CreatedAt *LocalTime      `gorm:"autoUpdateTime" json:"created_at"`
	UpdatedAt *LocalTime      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName FrameMaterial's table name
func (*FrameMaterial) TableName() string {
	return TableNameFrameMaterial
}
