package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(rg *gin.RouterGroup) {
	category := rg.Group("/categories")
	category.GET("/", controllers.GetPublishedCategories)
	category.GET("/options", controllers.GetCategoryOptions)
	category.Use(middleware.AdminRequireAuth(), middleware.RequireRole("admin"))
	{
		category.GET("/lists", controllers.GetCategories)
		category.GET("/:uuid", controllers.GetCategoryByUUID)
		category.PUT("/:uuid", controllers.EditCategory)
		category.POST("/submit", controllers.CreateCategory)
		category.DELETE("/:uuid", controllers.DeleteCategory)
		category.PATCH("/:uuid", controllers.DeleteImageCategory)
	}
}
