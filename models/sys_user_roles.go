// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameSysUserRole = "sys_user_roles"

// UserRole mapped from table <user_roles>
type SysUserRole struct {
	SysUserID int64 `gorm:"column:sys_user_id;primaryKey" json:"sys_user_id"`
	SysRoleID int64 `gorm:"column:sys_role_id;primaryKey" json:"sys_role_id"`
}

// TableName UserRole's table name
func (*SysUserRole) TableName() string {
	return TableNameSysUserRole
}
