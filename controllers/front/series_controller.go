package front

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProductInfo 产品信息结构体
type ProductInfo struct {
	ID            uint   `json:"id"`
	ModelNo       string `json:"model_no"`
	ImageURL      string `json:"image_url"`
	LensWidth     string `json:"lens_width"`
	NoseBridge    string `json:"nose_bridge"`
	TempleLength  string `json:"temple_length"`
	FrameMaterial string `json:"frame_material"`
}

// SeriesDetailResponse 系列详情响应结构体
type SeriesDetailResponse struct {
	SeriesName string        `json:"series_name"`
	Products   []ProductInfo `json:"products"`
}

// GetSeriesList 获取系列列表
func GetSeriesList(c *gin.Context) {
	// TODO: 实现获取系列列表逻辑
	// 1. 获取查询参数
	// 2. 查询系列列表
	c.JSON(http.StatusOK, gin.H{
		"message": "获取系列列表成功",
		// 添加系列列表数据
	})
}

// GetSeriesDetail 获取系列详情
func GetSeriesDetail(c *gin.Context) {
	seriesID := c.Param("id")

	// 使用工具函数检查用户是否为 VIP
	userID, _ := c.Get("userID")
	isVIP := utils.CheckUserRole(userID, "vip")

	// Query series with permission check
	var series models.Series
	query := config.DB.Where("id = ? AND deleted_at IS NULL", seriesID)
	if !isVIP {
		query = query.Where("is_new_design = 0 OR is_new_design IS NULL")
	}

	if err := query.First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, "Series not found or access denied", http.StatusNotFound)
			return
		}
		logrus.WithError(err).Error("Failed to query series details")
		utils.ErrorResponse(c, "Failed to get series details", http.StatusInternalServerError)
		return
	}

	// Query products with permission check
	var products []ProductInfo
	productsQuery := config.DB.Table("products").
		Select(`
			products.id,
			products.model_no,
			product_images.image_url,
			products.lens_width,
			products.nose_bridge,
			products.temple_length,
			frame_materials.name as frame_material
		`).
		Joins("LEFT JOIN frame_materials ON products.frame_material_id = frame_materials.id").
		Joins("LEFT JOIN product_images ON products.id = product_images.product_id").
		Where("products.series_id = ? AND products.deleted_at IS NULL", seriesID)

	err := productsQuery.Find(&products).Error
	if err != nil {
		logrus.WithError(err).Error("Failed to query product details")
		utils.ErrorResponse(c, "Failed to get product details", http.StatusInternalServerError)
		return
	}

	ossBaseURL := os.Getenv("OSS_BASE_URL")

	// 处理图片URL，添加OSS域名前缀
	for i := range products {
		if products[i].ImageURL != "" {
			products[i].ImageURL = fmt.Sprintf("%s/images/%s", ossBaseURL, products[i].ImageURL)
		}
	}

	response := SeriesDetailResponse{
		SeriesName: series.Name,
		Products:   products,
	}

	utils.SuccessResponse(c, "获取系列详情成功", response)
}
