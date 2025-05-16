package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func WebsiteRoutes(rg *gin.RouterGroup) {
	websites := rg.Group("/websites")

	websites.GET("/", controllers.GetUsers)

	// Route yang butuh admin
	adminOnly := websites.Group("/")
	adminOnly.Use(
		middleware.AdminRequireAuth(),
		middleware.RequireRole("admin"),
	)
	{
		adminOnly.POST("/submit", controllers.CreateUser)
		adminOnly.PUT("/:uuid", controllers.EditUser)
	}
}
