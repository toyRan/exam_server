package services

import (
	"exam_server/config"
	"exam_server/models"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// GetCategoryCascaderData 获取所有分类的级联选择器数据
func GetCategoryCascaderData() ([]map[string]interface{}, error) {
	var categories []models.Category

	// 查询顶级分类（PID 为 0 的分类）
	err := config.DB.Debug().Where("pid = 0").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	log.Printf("找到 %d 个顶级分类", len(categories))

	// 构造顶级分类的树结构
	cascaderData := []map[string]interface{}{}
	for _, category := range categories {
		cascaderData = append(cascaderData, buildCategoryTree(category, config.DB))
	}

	return cascaderData, nil
}

// 递归构造分类树
func buildCategoryTree(category models.Category, db *gorm.DB) map[string]interface{} {
	// 构造当前分类的数据
	node := map[string]interface{}{
		"value": category.ID,   // 分类 ID
		"label": category.Name, // 分类名称
	}

	// 查询子分类
	var children []models.Category
	result := db.Debug().Where("pid = ?", category.ID).Find(&children)
	if result.Error != nil {
		log.Printf("查询分类 %d 的子分类时出错: %v", category.ID, result.Error)
	}

	log.Printf("分类 %d (%s) 有 %d 个子分类", category.ID, category.Name, len(children))

	// 如果有子分类，递归构造
	if len(children) > 0 {
		childNodes := []map[string]interface{}{}
		for _, child := range children {
			childNodes = append(childNodes, buildCategoryTree(child, db))
		}
		node["children"] = childNodes
	}

	return node
}

// GetCategoryID 根据主分类和子分类名称获取对应的 category_id，如果没有找到则插入新记录
func GetCategoryID(mainCategoryName, subCategoryName string) (int64, error) {

	// 去除输入名称两边的空白
	mainCategoryName = strings.TrimSpace(mainCategoryName)
	subCategoryName = strings.TrimSpace(subCategoryName)

	var mainCategory models.Category
	var subCategory models.Category

	// 先查询主分类，如果不存在则创建
	if err := config.DB.Where("name = ? AND pid = 0", mainCategoryName).First(&mainCategory).Error; err != nil {
		// 如果没有找到主分类，创建新主分类
		mainCategory = models.Category{Name: mainCategoryName, Pid: 0} // PID 为 0
		if createErr := config.DB.Create(&mainCategory).Error; createErr != nil {
			log.Printf("failed to create main category '%s': %v", mainCategoryName, createErr)
			return 0, fmt.Errorf("failed to create main category '%s': %v", mainCategoryName, createErr)
		}
		fmt.Printf("Created new main category: %s with ID %d\n", mainCategoryName, mainCategory.ID)
	}

	// 查询子分类，如果不存在则创建
	if err := config.DB.Where("name = ? AND pid = ?", subCategoryName, mainCategory.ID).First(&subCategory).Error; err != nil {
		// 如果没有找到子分类，创建新子分类
		subCategory = models.Category{Name: subCategoryName, Pid: mainCategory.ID} // 子分类的 PID 为主分类的 ID
		if createErr := config.DB.Create(&subCategory).Error; createErr != nil {
			log.Printf("failed to create sub category '%s' under main category '%s': %v", subCategoryName, mainCategoryName, createErr)
			return 0, fmt.Errorf("failed to create sub category '%s' under main category '%s': %v", subCategoryName, mainCategoryName, createErr)
		}
		fmt.Printf("Created new sub category: %s under main category %s with ID %d\n", subCategoryName, mainCategoryName, subCategory.ID)
	}

	return subCategory.ID, nil
}

// GetCategoriesTree 获取完整分类树结构，并在内存中根据关键词过滤
func GetCategoriesTree(query string) ([]*models.Category, error) {
	var categories []models.Category

	// 一次性查询所有分类数据，不带任何关键词过滤
	if err := config.DB.Debug().Order("display_order").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 在内存中构建并过滤包含关键词的树结构
	return filterAndBuildCategoryTree(categories, query), nil
}

// filterAndBuildCategoryTree 构建并过滤符合条件的分类树结构
func filterAndBuildCategoryTree(categories []models.Category, query string) []*models.Category {
	categoryMap := make(map[int64]*models.Category)
	var rootCategories []*models.Category

	// 初始化 map 并设置空的 Children 切片
	for i := range categories {
		categories[i].Children = []*models.Category{}
		categoryMap[categories[i].ID] = &categories[i]
	}

	// 构建完整的父子关系树
	for i := range categories {
		category := categoryMap[categories[i].ID]
		if category.Pid == 0 {
			rootCategories = append(rootCategories, category)
		} else if parent, exists := categoryMap[category.Pid]; exists {
			parent.Children = append(parent.Children, category)
		}
	}

	// 过滤出包含关键词的树结构
	return filterCategories(rootCategories, query)
}

// filterCategories 递归过滤树，返回包含关键词的节点及其完整父子结构
func filterCategories(categories []*models.Category, query string) []*models.Category {
	var filteredCategories []*models.Category

	lowerQuery := strings.ToLower(query) // 将关键词转换为小写

	for _, category := range categories {
		// 递归过滤子节点
		category.Children = filterCategories(category.Children, query)

		// 如果当前节点包含关键词或有符合条件的子节点，则添加到结果中
		if strings.Contains(strings.ToLower(category.Name), lowerQuery) || len(category.Children) > 0 {
			filteredCategories = append(filteredCategories, category)
		}
	}

	return filteredCategories
}
