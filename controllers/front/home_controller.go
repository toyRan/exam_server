package front

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SeriesStatistics 系列统计信息
type SeriesStatistics struct {
	SeriesID     uint   `json:"series_id"`
	SeriesName   string `json:"series_name"`   // 系列名称
	MaterialName string `json:"material_name"` // 材质名称
	CategoryID   uint   `json:"category_id"`   // 类别ID
	CategoryName string `json:"category_name"` // 类别名称
	ModelsCount  int64  `json:"models_count"`  // 型号数量
	SkusCount    int64  `json:"skus_count"`    // SKU数量
	IsNewDesign  bool   `json:"is_new_design"` // 是否是最新设计
}

// CategoryWithCount 带数量的类别
type CategoryWithCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
	SeriesCount  int64  `json:"series_count"`
}

// MaterialWithCount 带数量的材质
type MaterialWithCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProductCount int64  `json:"product_count"`
	SeriesCount  int64  `json:"series_count"`
}

// GenderWithCount 带数量的性别选项
type GenderWithCount struct {
	Value        string `json:"value"`
	ProductCount int64  `json:"product_count"`
}

// FilterResponse 筛选条件响应结构
type FilterResponse struct {
	Categories     []CategoryWithCount `json:"categories"`
	FrameMaterials []MaterialWithCount `json:"frame_materials"`
}

// HomeResponse 首页响应数据结构
type HomeResponse struct {
	Filters       FilterResponse     `json:"filters"`
	OnlineCatalog []SeriesStatistics `json:"online_catalog"`
	Pagination    struct {
		Total    int64 `json:"total"`     // 总记录数
		Page     int   `json:"page"`      // 当前页码
		PageSize int   `json:"page_size"` // 每页数量
	} `json:"pagination"`
}

// HomeQueryParams 定义查询参数结构
type HomeQueryParams struct {
	CategoryID string `form:"category_id"`
	MaterialID string `form:"material_id"`
	Page       int    `form:"page,default=1"`       // 页码，默认1
	PageSize   int    `form:"page_size,default=10"` // 每页数量，默认10
}

// GetHomeData 获取首页数据
func GetHomeData(c *gin.Context) {
	var params HomeQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"params": c.Request.URL.Query(),
		}).Error("Failed to parse home page parameters")
		utils.ErrorResponse(c, "参数无效", http.StatusBadRequest)
	}

	// 获取用户VIP状态
	isVIP := false
	if userID, exists := c.Get("userID"); exists && userID != nil {
		isVIP = utils.CheckUserRole(userID, "vip")
	}

	var response HomeResponse
	db := config.DB

	// 实际的新品设计过滤逻辑已经在后面的 categoryQuery 和其他查询中实现

	// 1. 获取类别统计
	var categories []CategoryWithCount
	categoryQuery := db.Model(&models.Category{}).
		Select("categories.id, categories.name").
		Where("categories.deleted_at IS NULL")

	// 添加子查询统计产品数量
	productCountSubQuery := db.Model(&models.Product{}).
		Select("COUNT(DISTINCT products.id)").
		Joins("LEFT JOIN series ON products.series_id = series.id").
		Where("products.category_id = categories.id").
		Where("products.deleted_at IS NULL")
	if !isVIP {
		productCountSubQuery = productCountSubQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.MaterialID != "" {
		productCountSubQuery = productCountSubQuery.Where("products.frame_material_id = ?", params.MaterialID)
	}
	categoryQuery = categoryQuery.
		Select("categories.id, categories.name, (?) as product_count", productCountSubQuery)

	// 添加子查询统计系列数量
	seriesCountSubQuery := db.Model(&models.Series{}).
		Select("COUNT(DISTINCT series.id)").
		Joins("LEFT JOIN products ON series.id = products.series_id").
		Where("products.category_id = categories.id").
		Where("series.deleted_at IS NULL")
	if !isVIP {
		seriesCountSubQuery = seriesCountSubQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.MaterialID != "" {
		seriesCountSubQuery = seriesCountSubQuery.Where("products.frame_material_id = ?", params.MaterialID)
	}
	categoryQuery = categoryQuery.
		Select("categories.id, categories.name, (?) as product_count, (?) as series_count", productCountSubQuery, seriesCountSubQuery)

	if err := categoryQuery.Find(&categories).Error; err != nil {
		logrus.WithError(err).Error("Failed to get categories")
		utils.ErrorResponse(c, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 2. 获取材质统计
	var materials []MaterialWithCount
	materialQuery := db.Model(&models.FrameMaterial{}).
		Select("frame_materials.id, frame_materials.name").
		Where("frame_materials.deleted_at IS NULL")

	// 添加子查询统计产品数量
	productCountSubQuery = db.Model(&models.Product{}).
		Select("COUNT(DISTINCT products.id)").
		Joins("LEFT JOIN series ON products.series_id = series.id").
		Where("products.frame_material_id = frame_materials.id").
		Where("products.deleted_at IS NULL")
	if !isVIP {
		productCountSubQuery = productCountSubQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.CategoryID != "" {
		productCountSubQuery = productCountSubQuery.Where("products.category_id = ?", params.CategoryID)
	}
	materialQuery = materialQuery.
		Select("frame_materials.id, frame_materials.name, (?) as product_count", productCountSubQuery)

	// 添加子查询统计系列数量
	seriesCountSubQuery = db.Model(&models.Series{}).
		Select("COUNT(DISTINCT series.id)").
		Joins("LEFT JOIN products ON series.id = products.series_id").
		Where("products.frame_material_id = frame_materials.id").
		Where("series.deleted_at IS NULL")
	if !isVIP {
		seriesCountSubQuery = seriesCountSubQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.CategoryID != "" {
		seriesCountSubQuery = seriesCountSubQuery.Where("products.category_id = ?", params.CategoryID)
	}
	materialQuery = materialQuery.
		Select("frame_materials.id, frame_materials.name, (?) as product_count, (?) as series_count", productCountSubQuery, seriesCountSubQuery)

	if err := materialQuery.Find(&materials).Error; err != nil {
		logrus.WithError(err).Error("Failed to get materials")
		utils.ErrorResponse(c, "Failed to get materials", http.StatusInternalServerError)
		return
	}

	// 3. 获取系列统计数据
	var seriesStats []SeriesStatistics
	seriesQuery := db.Model(&models.Series{}).
		Select(`
			series.id as series_id,
			series.name as series_name,
			series.is_new_design,
			frame_materials.name as material_name,
			categories.id as category_id,
			categories.name as category_name,
			COUNT(DISTINCT products.id) as models_count,
			SUM(IFNULL(products.sku_count, 0)) as skus_count
		`).
		Joins("LEFT JOIN products ON series.id = products.series_id AND products.deleted_at IS NULL").
		Joins("LEFT JOIN frame_materials ON products.frame_material_id = frame_materials.id").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Where("series.deleted_at IS NULL").
		Group("series.id, series.name, frame_materials.name, categories.id, categories.name").
		Having("models_count > 0").
		Order("series.created_at DESC")

	if !isVIP {
		seriesQuery = seriesQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.CategoryID != "" {
		seriesQuery = seriesQuery.Where("products.category_id = ?", params.CategoryID)
	}
	if params.MaterialID != "" {
		seriesQuery = seriesQuery.Where("products.frame_material_id = ?", params.MaterialID)
	}

	// 获取总数
	var total int64
	countQuery := db.Model(&models.Series{}).
		Distinct("series.id").
		Joins("LEFT JOIN products ON series.id = products.series_id AND products.deleted_at IS NULL").
		Where("series.deleted_at IS NULL").
		Group("series.id").
		Having("COUNT(DISTINCT products.id) > 0")

	if !isVIP {
		countQuery = countQuery.Where("series.is_new_design = ? OR series.is_new_design IS NULL", false)
	}
	if params.CategoryID != "" {
		countQuery = countQuery.Where("products.category_id = ?", params.CategoryID)
	}
	if params.MaterialID != "" {
		countQuery = countQuery.Where("products.frame_material_id = ?", params.MaterialID)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		logrus.WithError(err).Error("Failed to get total count")
		utils.ErrorResponse(c, "获取总数失败", http.StatusInternalServerError)
	}

	// 设置筛选条件响应
	response.Filters = FilterResponse{
		Categories:     categories,
		FrameMaterials: materials,
	}

	// 添加分页
	offset := (params.Page - 1) * params.PageSize
	seriesQuery = seriesQuery.Limit(params.PageSize).Offset(offset)

	// 获取分页数据
	if err := seriesQuery.Find(&seriesStats).Error; err != nil {
		logrus.WithError(err).Error("Failed to get series list")
		utils.ErrorResponse(c, "获取系列列表失败", http.StatusInternalServerError)
		return
	}

	response.OnlineCatalog = seriesStats
	response.Pagination.Total = total
	response.Pagination.Page = params.Page
	response.Pagination.PageSize = params.PageSize

	utils.SuccessResponse(c, "获取数据成功", response)
}
