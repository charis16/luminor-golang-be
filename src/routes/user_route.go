package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/", controllers.GetUsers)
		users.POST("/", controllers.CreateUser)
	}
}
