package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/team-members", controllers.GetTeamMembers)
	users.GET("/options", controllers.GetUserOptions)
	users.Use(middleware.AdminRequireAuth(), middleware.RequireRole("admin"))
	{
		users.GET("/lists", controllers.GetUsers)
		users.GET("/:uuid", controllers.GetUserByUUID)
		users.PUT("/:uuid", controllers.EditUser)
		users.POST("/submit", controllers.CreateUser)
		users.DELETE("/:uuid", controllers.DeleteUser)
		users.PATCH("/:uuid", controllers.DeleteImageUser)
	}
}
