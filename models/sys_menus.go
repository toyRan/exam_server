package models

// SysMenu Menu 菜单模型
type SysMenu struct {
	ID       uint      `json:"id"`
	ParentID *uint     `json:"parent_id"` // 允许 ParentID 为 null
	Label    string    `json:"label"`
	LinkTo   string    `json:"link_to"`
	Icon     string    `json:"icon"`
	Order    int       `json:"order"`
	Children []SysMenu `json:"children" gorm:"-"` // 子菜单项，不存储在数据库中
}
