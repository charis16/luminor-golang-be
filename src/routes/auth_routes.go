package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/admin-login", controllers.AdminLogin)

		auth.POST("/admin-refresh-token", controllers.AdminRefreshToken)
		auth.POST("/admin-logout", controllers.AdminLogout)
		auth.POST("/admin-verify-token", controllers.AdminVerifyToken)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/admin-reset-password", controllers.AdminResetPassword)
	}
}
