// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameRolePermission = "role_permissions"

// RolePermission mapped from table <role_permissions>
type RolePermission struct {
	RoleID       int64 `gorm:"column:role_id;primaryKey" json:"role_id"`
	PermissionID int64 `gorm:"column:permission_id;primaryKey" json:"permission_id"`
}

// TableName RolePermission's table name
func (*RolePermission) TableName() string {
	return TableNameRolePermission
}
