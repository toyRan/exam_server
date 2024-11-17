package services

import (
	"exam_server/config"
	"exam_server/models"
	"strings"
)

// GetSysPermissionTree 获取完整权限树结构，并在内存中根据关键词过滤
func GetSysPermissionTree(query string) ([]*models.SysPermission, error) {
	var permissions []models.SysPermission

	// 一次性查询所有权限数据，不带任何关键词过滤
	if err := config.DB.Debug().Order("id").Find(&permissions).Error; err != nil {
		return nil, err
	}

	// 在内存中构建并过滤包含关键词的树结构
	return filterAndBuildPermissionTree(permissions, query), nil
}

// filterAndBuildPermissionTree 构建并过滤符合条件的权限树结构
func filterAndBuildPermissionTree(permissions []models.SysPermission, query string) []*models.SysPermission {
	permissionMap := make(map[int64]*models.SysPermission)
	var rootPermissions []*models.SysPermission

	// 初始化 map 并设置空的 Children 切片
	for i := range permissions {
		permissions[i].Children = []*models.SysPermission{}
		permissionMap[permissions[i].ID] = &permissions[i]
	}

	// 构建完整的父子关系树
	for i := range permissions {
		permission := permissionMap[permissions[i].ID]
		if permission.ParentID == 0 {
			rootPermissions = append(rootPermissions, permission)
		} else if parent, exists := permissionMap[permission.ParentID]; exists {
			parent.Children = append(parent.Children, permission)
		}
	}

	// 过滤出包含关键词的树结构
	return filterPermissions(rootPermissions, query)
}

// filterPermissions 递归过滤树，返回包含关键词的节点及其完整父子结构
func filterPermissions(permissions []*models.SysPermission, query string) []*models.SysPermission {
	var filteredPermissions []*models.SysPermission

	lowerQuery := strings.ToLower(query) // 将关键词转换为小写

	for _, permission := range permissions {
		// 递归过滤子节点
		permission.Children = filterPermissions(permission.Children, query)

		// 如果当前节点包含关键词或有符合条件的子节点，则添加到结果中
		if strings.Contains(strings.ToLower(permission.Name), lowerQuery) || len(permission.Children) > 0 {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}

	return filteredPermissions
}
