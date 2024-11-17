package routes

import (
	"exam_server/controllers/admin"
	"exam_server/controllers/front"
	"exam_server/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// 将根路由移到最前面
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// 后台
	//不需要 JWT 中间件保护的公共路由
	publicRoutes := r.Group("/api/V1/admin")
	{
		publicRoutes.POST("/login", admin.SysUserLogin)                          // 后台用户登录
		publicRoutes.GET("/captcha/refresh", admin.RefreshCaptcha)               // 刷新验证码
		publicRoutes.POST("/request-password-reset", admin.RequestPasswordReset) // 请求重置密码
		publicRoutes.POST("/reset-password", admin.ResetPassword)                // 重置密码
	}

	// 后台
	//需要 JWT 中间件保护的私有路由
	apiRoutes := r.Group("/api/V1/admin", middlewares.JWTAuthMiddleware())
	{
		//获取管理员菜单-仅自己权限内的
		apiRoutes.GET("/sys-menus", admin.GetSysUserMenus)

		//后台用户路由组
		sysUserRoutes := apiRoutes.Group("/sys-users")
		{
			sysUserRoutes.POST("/create", admin.SysUserStore)  // 创建后台用户
			sysUserRoutes.POST("/update", admin.SysUserUpdate) // 更新后台用户
			sysUserRoutes.POST("/delete", admin.SysUserDelete) // 删除后台用户
			sysUserRoutes.GET("", admin.GetAllSysUsers)        //所有后台用户 带分页
		}

		//后台角色路由组
		sysRoleRoutes := apiRoutes.Group("/sys-roles")
		{
			sysRoleRoutes.POST("/create", admin.SysRoleCreate)                  // 创建管理员角色
			sysRoleRoutes.POST("/delete", admin.SysRoleDelete)                  // 删除管理员角色
			sysRoleRoutes.POST("/update", admin.SysRoleUpdate)                  // 更新管理员角色
			sysRoleRoutes.GET("", admin.GetAllSysRoles)                         // 查询所有角色（无分页）
			sysRoleRoutes.GET("/paginated", admin.GetAllSysRolesPaginated)      // 查询所有角色（分页）
			sysRoleRoutes.GET("/:id", admin.GetSysRole)                         // 查询单个角色
			sysRoleRoutes.POST("/set-permissions", admin.SetSysRolePermissions) // 设置角色权限
		}

		//后台权限
		sysPermissionRoutes := apiRoutes.Group("/sys-permissions")
		{
			sysPermissionRoutes.GET("/tree", admin.GetSysPermissionsTree)              // 查询所有权限（分页）
			sysPermissionRoutes.POST("/create", admin.SysPermissionCreate)             // 创建权限
			sysPermissionRoutes.POST("/delete", admin.SysPermissionDelete)             // 删除权限
			sysPermissionRoutes.POST("/update", admin.SysPermissionUpdate)             // 更新权限
			sysPermissionRoutes.GET("", admin.GetAllSysPermissions)                    // 查询所有权限（注意这里与下一行顺序）
			sysPermissionRoutes.GET("/paginated", admin.GetAllSysPermissionsPaginated) // 查询所有权限（分页）
			sysPermissionRoutes.GET("/:id", admin.GetSysPermission)                    // 查询单个权限

		}

		// 商品分类
		categoryRoutes := apiRoutes.Group("/categories")
		{
			categoryRoutes.GET("/fetch-cascade", admin.GetCategoriesForCascader) //分类级联 （下拉绑定用）
			categoryRoutes.GET("/tree", admin.GetCategoriesTree)                 //分类列表 （后台管理用）
			categoryRoutes.POST("/create", admin.CreateCategory)                 // 创建新分类
			categoryRoutes.POST("/update", admin.UpdateCategory)                 // 更新分类
			categoryRoutes.POST("/delete", admin.DeleteCategory)                 // 删除分类
			categoryRoutes.GET("", admin.GetCategories)                          // 获取所有分类
			categoryRoutes.GET("/:id", admin.GetCategory)                        // 查询单个分类
		}

		// 商品材质
		materialRoutes := apiRoutes.Group("/frame-materials")
		{
			materialRoutes.POST("/create", admin.CreateFrameMaterial)            // 创建新框材质
			materialRoutes.POST("/update", admin.UpdateFrameMaterial)            // 更新框材质
			materialRoutes.POST("/delete", admin.DeleteFrameMaterial)            // 删除框材质
			materialRoutes.GET("", admin.GetAllFrameMaterials)                   // 获取所有框材质
			materialRoutes.GET("paginated", admin.GetAllFrameMaterialsPaginated) // 获取所有框材质
			materialRoutes.GET("/:id", admin.GetFrameMaterial)                   // 查询单个框材质
		}

		// 商品系列
		seriesRoutes := apiRoutes.Group("/series")
		{
			seriesRoutes.POST("/create", admin.CreateSeries)
			seriesRoutes.POST("/update", admin.UpdateSeries)
			seriesRoutes.POST("/delete", admin.DeleteSeries)
			seriesRoutes.POST("/delete-batch", admin.DeleteSeriesBatch)
			seriesRoutes.GET("", admin.GetAllSeries)
			seriesRoutes.GET("/paginated", admin.GetAllSeriesPaginated)
			seriesRoutes.GET("/:id", admin.GetSeries)
		}

		// 商品品牌
		brandRoutes := apiRoutes.Group("/brands")
		{
			brandRoutes.POST("/create", admin.CreateBrand)
			brandRoutes.POST("/update", admin.UpdateBrand)
			brandRoutes.POST("/delete", admin.DeleteBrand)
			brandRoutes.POST("/delete-batch", admin.DeleteBrandsBatch)
			brandRoutes.GET("", admin.GetAllBrands)
			brandRoutes.GET("/paginated", admin.GetAllBrandsPaginated)
			brandRoutes.GET("/:id", admin.GetBrand)
		}

		// 产品
		productRoutes := apiRoutes.Group("/products")
		{
			productRoutes.POST("/create", admin.CreateProduct)
			productRoutes.POST("/update", admin.UpdateProduct)
			productRoutes.POST("/delete", admin.DeleteProduct)
			productRoutes.POST("/delete-batch", admin.DeleteProductBatch)
			productRoutes.GET("/paginated", admin.GetAllProductsPaginated)
			productRoutes.GET("", admin.GetAllProductsPaginated) //获取全部products数据
			productRoutes.GET("/:id", admin.GetProduct)          //获取单个商品数据
		}

		// 上传文件到OSS
		ossRoutes := apiRoutes.Group("/oss")
		{
			ossRoutes.POST("/upload/multiple", admin.UploadFiles)
		}

		//获取所有前台用户列表
		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.POST("/create", admin.CreateUser) // 创建用户
			userRoutes.POST("/update", admin.UpdateUser) // 更新用户
			userRoutes.POST("/delete", admin.DeleteUser) // 删除用户
			//userRoutes.POST("/delete-batch", admin.DeleteUserBatch)
			userRoutes.GET("", admin.GetUsersPaginated) // 获取所有前台用户列表
			//userRoutes.GET("/:id", admin.GetUser)          // 获取所有前台用户列表
		}

		//获取所有前台角色列表
		roleRoutes := apiRoutes.Group("/roles")
		{
			roleRoutes.POST("/create", admin.CreateRole) // 创建前台角色
			roleRoutes.POST("/update", admin.UpdateRole) // 更新前台角色
			roleRoutes.POST("/delete", admin.DeleteRole) // 删除前台角色
			//roleRoutes.POST("/delete-batch", admin.DeleteRoleBatch)
			roleRoutes.GET("/paginated", admin.GetRolesPaginated) // 获取所有前台角色
			roleRoutes.GET("/all", admin.GetRoleList)             // 获取所有前台角色 不分页
			roleRoutes.GET("/:id", admin.GetRole)                 // 获取单个前台角色信息
		}

	}

	// 前台
	// 不需要 JWT 中间件保护的公共路由
	frontPublicRoutes := r.Group("/api/V1")
	{
		// 首页相关
		frontPublicRoutes.GET("/home", middlewares.OptionalJWTAuth(), front.GetHomeData)

		// 用户认证相关
		frontPublicRoutes.POST("/register", front.UserRegister) // 用户注册
		frontPublicRoutes.POST("/login", front.UserLogin)       // 用户登录

		// 商品相关
		// frontPublicRoutes.GET("/products", front.GetProducts)          // 获取商品列表
		// frontPublicRoutes.GET("/products/:id", front.GetProductDetail) // 获取商品详情
	}

	// 需要 JWT 中间件保护的私有路由
	frontPrivateRoutes := r.Group("/api/V1", middlewares.JWTAuthMiddleware())
	{
		// 用户相关
		userRoutes := frontPrivateRoutes.Group("/user")
		{
			userRoutes.GET("/profile", front.GetUserProfile)        // 获取用户信息
			userRoutes.POST("/profile/update", front.UpdateProfile) // 更新用户信息
		}

		// frontPublicRoutes.POST("/logout", front.UserLogout) // 用户退出

		// PDF下载相关
		// pdfRoutes := frontPrivateRoutes.Group("/pdf")
		// {
		// 	pdfRoutes.GET("/download/:id", front.DownloadPDF)   // 下载PDF
		// 	pdfRoutes.GET("/history", front.GetDownloadHistory) // 获取下载历史
		// }
	}

}
