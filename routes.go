package main

import (
	"github.com/gin-gonic/gin"
	"ngxs.site/ginessential/controller"
	"ngxs.site/ginessential/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine{
	r.Use(middleware.CORSMiddleware(),middleware.RecoveryMiddleware())
	r.POST("/api/auth/register",controller.Register)
	r.POST("/api/auth/login",controller.Login)
	r.GET("/api/auth/info",middleware.AuthMiddleware(),controller.Info)


	categoryRoutes := r.Group("/categories")
	categoryController := controller.NewCategoryController()
	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update)  // 跟patch的区别就是，update是替换，patch是部分修改
	categoryRoutes.GET("/:id", categoryController.Show)
	categoryRoutes.DELETE("/:id", categoryController.Delete)
	//categoryRoutes.PATCH("/:id", )  // 部分修改

	postRoutes := r.Group("/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	postController := controller.NewPostController()
	postRoutes.POST("", postController.Create)
	postRoutes.PUT("/:id", postController.Update)  // 跟patch的区别就是，update是替换，patch是部分修改
	postRoutes.GET("/:id", postController.Show)
	postRoutes.DELETE("/:id", postController.Delete)
	postRoutes.POST("page/list",postController.PageList)

	fileRoutes := r.Group("/files")
	fileRoutes.Use(middleware.AuthMiddleware())
	fileController := controller.NewFileController()
	fileRoutes.POST("",fileController.Create)  //上传文件

	return r
}