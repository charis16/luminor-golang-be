package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func AlbumRoutes(rg *gin.RouterGroup) {
	faq := rg.Group("/albums")
	faq.Use(middleware.AdminRequireAuth(), middleware.RequireRole("admin"))
	{
		faq.GET("/lists", controllers.GetAlbums)
		faq.GET("/:uuid", controllers.GetAlbumByUUID)
		faq.PUT("/:uuid", controllers.EditAlbum)
		faq.POST("/submit", controllers.CreateAlbum)
		faq.DELETE("/:uuid", controllers.DeleteAlbum)
	}
}
