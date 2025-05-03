package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.POST("/refresh-token", controllers.RefreshToken)
		auth.POST("/logout", controllers.Logout)
		auth.POST("/verify-token", controllers.VerifyToken)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
	}
}
