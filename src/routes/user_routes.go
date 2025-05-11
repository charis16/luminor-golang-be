package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.Use(middleware.RequireAuth(), middleware.RequireRole("admin"))
	{
		users.GET("/lists", controllers.GetUsers)
		users.GET("/:uuid", controllers.GetUserByUUID)
		users.PUT("/edit/:uuid", controllers.EditUser)
		users.GET("image", controllers.ProxyUserImage)
		users.POST("/submit", controllers.CreateUser)
		users.DELETE("/:uuid", controllers.DeleteUser)
	}
}
