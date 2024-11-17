package models

import "github.com/guregu/null/v5"

const TableNameProductImage = "product_images"

// ProductImage mapped from table <product_skus>
type ProductImage struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ProductID null.Int64 `gorm:"column:product_id" json:"product_id"`
	ImageURL  string     `gorm:"column:image_url" json:"image_url"`
	CreatedAt *LocalTime `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *LocalTime `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName ProductSku's table name
func (*ProductImage) TableName() string {
	return TableNameProductImage
}
