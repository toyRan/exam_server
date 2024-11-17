package services

import (
	"exam_server/config"
	"exam_server/models"
)

// GetUserSysMenus 获取用户的菜单权限
func GetUserSysMenus(userId int64) (*[]*models.SysMenu, error) {
	var menus []models.SysMenu

	// 执行原始 SQL 查询获取菜单项
	err := config.DB.Raw(`
       SELECT DISTINCT m.id, m.parent_id, m.label, m.link_to, m.icon, m.order
      FROM sys_menus m
      INNER JOIN sys_role_menu rm ON m.id = rm.sys_menu_id
      INNER JOIN sys_roles r ON rm.sys_role_id = r.id
      INNER JOIN sys_user_role ur ON r.id = ur.sys_role_id
      WHERE ur.sys_user_id = ?
      AND m.deleted_at IS NULL
      ORDER BY m.order DESC
    `, userId).Scan(&menus).Error

	if err != nil {
		return nil, err
	}

	//// 在打印菜单之前格式化输出
	//menuData, err := json.MarshalIndent(menus, "", "  ")
	//if err != nil {
	//	log.Printf("Error formatting menus for output: %v", err)
	//} else {
	//	log.Printf("GetUserSysMenus:\n%s", menuData)
	//}

	// 构建菜单树
	menuTree := BuildSysMenuTree(menus)
	return &menuTree, nil
}

// BuildSysMenuTree 构建菜单树
func BuildSysMenuTree(menus []models.SysMenu) []*models.SysMenu {
	menuMap := make(map[uint]*models.SysMenu) // 用于快速查找父菜单项
	var menuTree []*models.SysMenu            // 用于存储最终的树结构（使用指针切片）

	// 将菜单项映射到 map 中并初始化 Children
	for i := range menus {
		menus[i].Children = []models.SysMenu{} // 初始化 Children 字段为空切片
		menuMap[menus[i].ID] = &menus[i]       // 将菜单项加入 map 中，键是菜单项的 ID
	}

	// 构建树结构
	for i := range menus {
		currentSysMenu := menuMap[menus[i].ID]
		if currentSysMenu.ParentID == nil || *currentSysMenu.ParentID == 0 {
			// 如果没有父级，则是顶级菜单，直接添加到 menuTree 中
			menuTree = append(menuTree, currentSysMenu)
		} else {
			// 如果有父级，则查找父菜单，将当前菜单加入到父菜单的 Children 中
			if parentSysMenu, exists := menuMap[*currentSysMenu.ParentID]; exists {
				parentSysMenu.Children = append(parentSysMenu.Children, *currentSysMenu)
			}
		}
	}

	return menuTree // 返回构建好的完整菜单树
}
