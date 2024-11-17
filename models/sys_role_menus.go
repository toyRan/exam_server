package models

const TableNameSysRoleMenu = "sys_role_menus"

// SysRoleMenu 中间表
type SysRoleMenu struct {
	SysRoleID uint `json:"sys_role_id"`
	SysMenuID uint `json:"sys_menu_id"`
}

// TableName User's table name
func (*SysRoleMenu) TableName() string {
	return TableNameSysRoleMenu
}
