package admin

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/services"
	"exam_server/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"
)

// ProductResponse 定义扁平化的响应结构
type ProductResponse struct {
	ID           int64       `json:"id"`
	ModelNo      string      `json:"model_no"`
	ItemCode     null.String `json:"item_code"`
	Title        string      `json:"title"`
	Gender       null.String `json:"gender"`
	LensWidth    float32     `json:"lens_width"`
	NoseBridge   float32     `json:"nose_bridge"`
	TempleLength float32     `json:"temple_length"`

	SeriesID          null.Int64 `json:"series_id"`
	SkuCount          int64      `json:"sku_count"`
	SeriesName        string     `json:"series_name"`
	CategoryID        int64      `json:"category_id"`
	CategoryName      string     `json:"category_name"`
	BrandID           null.Int64 `json:"brand_id"`
	BrandName         string     `json:"brand_name"`
	FrameMaterialID   int64      `json:"frame_material_id"`
	FrameMaterialName string     `json:"frame_material_name"`
	CategoryPath      []int      `json:"category_path"`
	Description       string     `json:"description"`

	ImageURLs []string          `json:"image_urls"`
	CreatedAt *models.LocalTime `json:"created_at"`
	UpdatedAt *models.LocalTime `json:"updated_at"`
}

// GetAllProductsPaginated 获取所有商品（支持分页）
func GetAllProductsPaginated(c *gin.Context) {
	// 获取查询参数 page 和 pageSize
	queryStr := c.DefaultQuery("query", "")
	pageStr := c.DefaultQuery("currentPage", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 分页查询
	var products []models.Product
	var total int64

	// 构建查询
	query := config.DB.Model(&models.Product{})

	// 修改查询以预加载关联数据
	query = query.
		Preload("Series").
		Preload("Category").
		Preload("Brand").
		Preload("FrameMaterial")

	// 如果 query 参数存在，则根据 ItemCode 或 Title 进行模糊搜索
	if queryStr != "" {
		searchPattern := queryStr + "%"
		query = query.Where("item_code LIKE ? OR model_no LIKE ?", searchPattern, searchPattern)
	}

	// 获取总数
	query.Count(&total)

	// 进行分页
	offset := (page - 1) * pageSize
	//config.DB.Limit(pageSize).Offset(offset).Find(&products)

	query.Limit(pageSize).Offset(offset).Find(&products)

	// 从环境变量获取域名
	ossBaseURL := os.Getenv("OSS_BASE_URL") // 例如 "https://jx-ebook.oss-cn-hongkong.aliyuncs.com"

	// 遍历每个产品，查找其主图
	for i, product := range products {
		var mainFile models.ProductImage

		// 获取主图
		if err := config.DB.Where("product_id = ?", product.ID).First(&mainFile).Error; err != nil {
			log.Println("Failed to fetch main file for product", product.ID, err)
			continue
		}

		// 构建完整的URL并附加到产品数据，添加 images/ 前缀
		fullImageURL := fmt.Sprintf("%s/images/%s", ossBaseURL, mainFile.ImageURL)
		products[i].ImageURLs = append(products[i].ImageURLs, fullImageURL)
	}

	// 定义扁平化的响应结构
	type ProductResponse struct {
		ID           int64       `json:"id"`
		ModelNo      string      `json:"model_no"`
		ItemCode     null.String `json:"item_code"`
		Title        string      `json:"title"`
		Gender       null.String `json:"gender"`
		LensWidth    float32     `json:"lens_width"`
		NoseBridge   float32     `json:"nose_bridge"`
		TempleLength float32     `json:"temple_length"`

		SeriesID          null.Int64 `json:"series_id"`
		SkuCount          int64      `json:"sku_count"`
		SeriesName        string     `json:"series_name"`
		CategoryID        int64      `json:"category_id"`
		CategoryName      string     `json:"category_name"`
		BrandID           null.Int64 `json:"brand_id"`
		BrandName         string     `json:"brand_name"`
		FrameMaterialID   int64      `json:"frame_material_id"`
		FrameMaterialName string     `json:"frame_material_name"`

		ImageURLs []string          `json:"image_urls"`
		CreatedAt *models.LocalTime `json:"created_at"`
		UpdatedAt *models.LocalTime `json:"updated_at"`
	}

	// 初始化空数组，确保即使没有数据也会返回空数组而不是 null
	responseProducts := make([]ProductResponse, 0)

	// 转换为响应格式
	for _, p := range products {
		response := ProductResponse{
			ID:           p.ID,
			ModelNo:      p.ModelNO,
			ItemCode:     p.ItemCode,
			Title:        p.Title,
			Gender:       p.Gender,
			LensWidth:    p.LensWidth,
			NoseBridge:   p.NoseBridge,
			TempleLength: p.TempleLength,

			SeriesID:          p.SeriesID,
			SkuCount:          p.SkuCount,
			SeriesName:        getSafeString(p.Series, "Name"),
			CategoryID:        p.CategoryID,
			CategoryName:      getSafeString(p.Category, "Name"),
			BrandID:           p.BrandID,
			BrandName:         getSafeString(p.Brand, "Name"),
			FrameMaterialID:   p.FrameMaterialID,
			FrameMaterialName: getSafeString(p.FrameMaterial, "Name"),

			ImageURLs: p.ImageURLs,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
		responseProducts = append(responseProducts, response)
	}

	// 构造响应
	data := gin.H{
		"products":    responseProducts,
		"total":       total,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalPage":   (total + int64(pageSize) - 1) / int64(pageSize),
	}

	// 返回成功响应
	utils.SuccessResponse(c, "Products fetched successfully", data)
}

// getSafeString 安全地获取关联对象的字符串属性
func getSafeString(obj interface{}, field string) string {
	if obj == nil {
		return ""
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ""
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return ""
	}

	return fieldVal.String()
}

// GetProduct 根据 ID 获取单个商品
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.
		Preload("Series").
		Preload("Category").
		Preload("Brand").
		Preload("FrameMaterial").
		First(&product, id).Error; err != nil {
		utils.ErrorResponse(c, "Product not found", http.StatusNotFound)
		return
	}

	// 获取分类路径
	var categoryPath []int
	currentCategoryID := product.CategoryID
	for currentCategoryID > 0 {
		var category models.Category
		if err := config.DB.First(&category, currentCategoryID).Error; err != nil {
			break
		}
		// 将当前分类ID添加到路径开头
		categoryPath = append([]int{int(currentCategoryID)}, categoryPath...)
		currentCategoryID = category.Pid
	}

	// 查询所有相关文件
	var productFiles []models.ProductImage
	if err := config.DB.Where("product_id = ?", product.ID).Find(&productFiles).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch product files", http.StatusInternalServerError)
		return
	}

	// 从环境变量获取域名
	ossBaseURL := os.Getenv("OSS_BASE_URL")

	// 构建完整的图片URL列表
	var imageURLs []string
	for _, file := range productFiles {
		fullImageURL := fmt.Sprintf("%s/images/%s", ossBaseURL, file.ImageURL)
		imageURLs = append(imageURLs, fullImageURL)
	}

	// 将图片URL列表和分类路径添加到product对象中
	product.ImageURLs = imageURLs
	product.CategoryPath = categoryPath

	// 转换为扁平化响应
	response := ProductResponse{
		ID:           product.ID,
		ModelNo:      product.ModelNO,
		ItemCode:     product.ItemCode,
		Title:        product.Title,
		Gender:       product.Gender,
		LensWidth:    product.LensWidth,
		NoseBridge:   product.NoseBridge,
		TempleLength: product.TempleLength,

		SeriesID:          product.SeriesID,
		SkuCount:          product.SkuCount,
		SeriesName:        getSafeString(product.Series, "Name"),
		CategoryID:        product.CategoryID,
		CategoryName:      getSafeString(product.Category, "Name"),
		BrandID:           product.BrandID,
		BrandName:         getSafeString(product.Brand, "Name"),
		FrameMaterialID:   product.FrameMaterialID,
		FrameMaterialName: getSafeString(product.FrameMaterial, "Name"),
		Description:       product.Description,

		ImageURLs:    imageURLs,
		CategoryPath: categoryPath,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	// 直接返回product对象
	utils.SuccessResponse(c, "Product fetched successfully", gin.H{"product": response})
}

// CreateProduct 创建商品
func CreateProduct(c *gin.Context) {
	var request struct {
		models.Product
		ImageURLs []string `json:"image_urls" binding:"required,min=1"`
	}

	// 参数验证
	if err := c.ShouldBindJSON(&request); err != nil {
		validationErrors := utils.HandleValidationError(err)
		if len(validationErrors) > 0 {
			c.JSON(400, gin.H{
				"success": false,
				"message": "参数验证失败",
				"errors":  validationErrors,
			})
			return
		}
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 验证图片
	validImageUrls := make([]string, 0)
	for _, url := range request.ImageURLs {
		if url != "" {
			validImageUrls = append(validImageUrls, url)
		}
	}

	if len(validImageUrls) == 0 {
		utils.ErrorResponse(c, "At least one product image is required", http.StatusBadRequest)
		return
	}

	// 如果指定了系列ID，验证框材质是否匹配
	if !request.Product.SeriesID.IsZero() {
		var series models.Series
		if err := config.DB.First(&series, request.Product.SeriesID).Error; err != nil {
			utils.ErrorResponse(c, "Series not found", http.StatusBadRequest)
			return
		}

		if series.FrameMaterialID != request.Product.FrameMaterialID {
			utils.ErrorResponse(c, "Frame material must match the series frame material", http.StatusBadRequest)
			return
		}
	}

	// 开始事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建产品
	if err := tx.Create(&request.Product).Error; err != nil {
		tx.Rollback()
		// 处理数据库错误
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	// 处理图片
	for _, imageURL := range request.ImageURLs {
		filename := filepath.Base(imageURL)
		productImage := models.ProductImage{
			ProductID: null.IntFrom(int64(request.Product.ID)),
			ImageURL:  filename,
		}

		if err := tx.Create(&productImage).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to create product image: %v", err)
			utils.ErrorResponse(c, "Failed to create product images", http.StatusInternalServerError)
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	// 重新查询完整的产品信息，包括关联数据
	var product models.Product
	if err := config.DB.
		Preload("Series").
		Preload("Category").
		Preload("Brand").
		Preload("FrameMaterial").
		First(&product, request.Product.ID).Error; err != nil {
		log.Println("查看关联数据失败了")
		utils.ErrorResponse(c, "Failed to fetch created product", http.StatusInternalServerError)
		return
	}

	// 构建图片URL列表
	ossBaseURL := os.Getenv("OSS_BASE_URL")
	var imageURLs []string
	for _, url := range request.ImageURLs {
		filename := filepath.Base(url)
		fullImageURL := fmt.Sprintf("%s/images/%s", ossBaseURL, filename)
		imageURLs = append(imageURLs, fullImageURL)
	}
	log.Printf("imageURLs: %v", imageURLs)

	// 转换为 ProductResponse 格式
	//response := ProductResponse{
	//	ID:           product.ID,
	//	ModelNo:      product.ModelNO,
	//	ItemCode:     product.ItemCode,
	//	Title:        product.Title,
	//	Gender:       product.Gender,
	//	LensWidth:    product.LensWidth,
	//	NoseBridge:   product.NoseBridge,
	//	TempleLength: product.TempleLength,
	//
	//	SeriesID:          product.SeriesID,
	//	SkuCount:          product.SkuCount,
	//	SeriesName:        product.Series.Name,
	//	CategoryID:        product.CategoryID,
	//	CategoryName:      product.Category.Name,
	//	BrandID:           product.BrandID,
	//	BrandName:         product.Brand.Name,
	//	FrameMaterialID:   product.FrameMaterialID,
	//	FrameMaterialName: product.FrameMaterial.Name,
	//
	//	ImageURLs: imageURLs,
	//	CreatedAt: product.CreatedAt,
	//	UpdatedAt: product.UpdatedAt,
	//}
	// 添加更多日志来定位问题
	log.Printf("开始构建响应")

	// 使用安全的方式构建响应
	response := ProductResponse{
		ID:           product.ID,
		ModelNo:      product.ModelNO,
		ItemCode:     product.ItemCode,
		Title:        product.Title,
		Gender:       product.Gender,
		LensWidth:    product.LensWidth,
		NoseBridge:   product.NoseBridge,
		TempleLength: product.TempleLength,
	}

	// 分步设置关联字段，并添加日志
	log.Printf("设置基本字段完成")

	// 安全设置 Series 相关字段
	if !product.SeriesID.IsZero() {
		response.SeriesID = product.SeriesID
		if product.Series != nil {
			response.SeriesName = product.Series.Name
		}
	}

	log.Printf("设置 Series 字段完成")

	// 设置其他关联字段
	response.CategoryID = product.CategoryID
	if product.Category != nil {
		response.CategoryName = product.Category.Name
	}

	log.Printf("设置 Category 字段完成")

	if !product.BrandID.IsZero() {
		response.BrandID = product.BrandID
		if product.Brand != nil {
			response.BrandName = product.Brand.Name
		}
	}

	log.Printf("设置 Brand 字段完成")

	response.FrameMaterialID = product.FrameMaterialID
	if product.FrameMaterial != nil {
		response.FrameMaterialName = product.FrameMaterial.Name
	}

	log.Printf("设置 FrameMaterial 字段完成")

	// 设置图片URLs和时间戳
	response.ImageURLs = imageURLs
	response.CreatedAt = product.CreatedAt
	response.UpdatedAt = product.UpdatedAt

	log.Printf("设置响应完成")

	utils.SuccessResponse(c, "商品创建成功", response)

	log.Printf("响应已发送")
}

// UpdateProduct 更新商品
func UpdateProduct(c *gin.Context) {
	var request struct {
		models.Product
		ImageURLs        []string `json:"image_urls"`
		DeletedImageURLs []string `json:"deleted_image_urls"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 如果指定了系列ID，验证框材质是否匹配
	if !request.Product.SeriesID.IsZero() {
		var series models.Series
		if err := config.DB.First(&series, request.Product.SeriesID).Error; err != nil {
			utils.ErrorResponse(c, "框材质不存在", http.StatusBadRequest)
			return
		}

		if series.FrameMaterialID != request.Product.FrameMaterialID {
			utils.ErrorResponse(c, "产品框材质必须与系列框材质一致", http.StatusBadRequest)
			return
		}
	}

	// 开始事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新产品基本信息
	if err := tx.Model(&models.Product{}).Where("id = ?", request.ID).Updates(request.Product).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "更新产品信息失败", http.StatusInternalServerError)
		return
	}

	// 2. 处理要删除的图片
	if len(request.DeletedImageURLs) > 0 {
		// 从 OSS 删除文件
		if err := services.DeleteFilesFromOSS(request.DeletedImageURLs); err != nil {
			log.Printf("Warning: Failed to delete files from OSS: %v", err)
		}

		// 从数据库删除记录
		deleteFilenames := make([]string, 0, len(request.DeletedImageURLs))
		for _, url := range request.DeletedImageURLs {
			filename := filepath.Base(url)
			deleteFilenames = append(deleteFilenames, filename)
		}

		if err := tx.Where("product_id = ? AND image_url IN ?", request.ID, deleteFilenames).
			Delete(&models.ProductImage{}).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "Failed to delete image records", http.StatusInternalServerError)
			return
		}
	}

	// 3. 处理图片记录
	// 3.1 获取现有图片记录
	var existingImages []models.ProductImage
	if err := tx.Where("product_id = ?", request.ID).Find(&existingImages).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to fetch existing images", http.StatusInternalServerError)
		return
	}

	// 3.2 创建现有图片map，用于检查重复
	existingImageMap := make(map[string]bool)
	for _, img := range existingImages {
		existingImageMap[img.ImageURL] = true
	}

	// 3.3 只添加新的图片记录
	for _, imageURL := range request.ImageURLs {
		filename := filepath.Base(imageURL)
		// 如果图片已存在，跳过
		if existingImageMap[filename] {
			continue
		}

		newImage := models.ProductImage{
			ProductID: null.IntFrom(request.ID),
			ImageURL:  filename,
		}
		if err := tx.Create(&newImage).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "Failed to create new image record", http.StatusInternalServerError)
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// 重新加载产品信息
	var product models.Product
	config.DB.
		Preload("Series").
		Preload("Category").
		Preload("Brand").
		Preload("FrameMaterial").
		First(&product, request.ID)

	// 构建响应
	response := ProductResponse{
		ID:           product.ID,
		ModelNo:      product.ModelNO,
		ItemCode:     product.ItemCode,
		Title:        product.Title,
		Gender:       product.Gender,
		LensWidth:    product.LensWidth,
		NoseBridge:   product.NoseBridge,
		TempleLength: product.TempleLength,

		SeriesID:          product.SeriesID,
		SkuCount:          product.SkuCount,
		SeriesName:        getSafeString(product.Series, "Name"),
		CategoryID:        product.CategoryID,
		CategoryName:      getSafeString(product.Category, "Name"),
		BrandID:           product.BrandID,
		BrandName:         getSafeString(product.Brand, "Name"),
		FrameMaterialID:   product.FrameMaterialID,
		FrameMaterialName: getSafeString(product.FrameMaterial, "Name"),

		ImageURLs: request.ImageURLs,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	utils.SuccessResponse(c, "Product updated successfully", response)
}

// DeleteProduct 删除商品
func DeleteProduct(c *gin.Context) {
	var request struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. 先查询所有相关的图片记录
	var productImages []models.ProductImage
	if err := config.DB.Where("product_id = ?", request.ID).Find(&productImages).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch product images", http.StatusInternalServerError)
		return
	}

	// 2. 开始事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2.1 删除图片记录
	if err := tx.Where("product_id = ?", request.ID).Delete(&models.ProductImage{}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to delete product images", http.StatusInternalServerError)
		return
	}

	// 2.2 删除产品记录
	if err := tx.Delete(&models.Product{}, request.ID).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// 2.3 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// 3. 删除 OSS 文件
	if len(productImages) > 0 {
		ossUrls := make([]string, 0, len(productImages))
		ossBaseURL := os.Getenv("OSS_BASE_URL")
		for _, img := range productImages {
			ossUrls = append(ossUrls, fmt.Sprintf("%s/images/%s", ossBaseURL, img.ImageURL))
		}

		if err := services.DeleteFilesFromOSS(ossUrls); err != nil {
			// 记录错误但不影响返回结果，因为数据库记录已经删除
			log.Printf("Warning: Failed to delete files from OSS: %v", err)
		}
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

// DeleteProductBatch 批量删除
func DeleteProductBatch(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"` // 从请求体中获取 ID 列表
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 执行软删除
	if err := config.DB.Where("id IN ?", request.IDs).Delete(&models.Product{}).Error; err != nil {
		log.Printf("Failed to delete products: %v", err)
		utils.ErrorResponse(c, "Failed to delete products", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "商品信息批量删除成功", nil)
}
