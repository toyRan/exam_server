// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameProductMaterial = "product_material"

// ProductMaterial mapped from table <product_material>
type ProductMaterial struct {
	ProductID       int64 `gorm:"column:product_id;primaryKey" json:"product_id"`
	FrameMaterialID int64 `gorm:"column:frame_material_id;primaryKey" json:"frame_material_id"`
}

// TableName ProductMaterial's table name
func (*ProductMaterial) TableName() string {
	return TableNameProductMaterial
}
