package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func AlbumRoutes(rg *gin.RouterGroup) {
	albums := rg.Group("/albums")
	albums.GET("/", controllers.GetLatestAlbum)
	albums.GET("/category/:slug", controllers.GetAlbumByCategorySlug)
	// albums.GET("/portfolio/:slug", controllers.GetAlbumByPortfolioSlug)
	albums.Use(middleware.AdminRequireAuth(), middleware.RequireRole("admin"))
	{
		albums.GET("/lists", controllers.GetAlbums)
		albums.GET("/:uuid", controllers.GetAlbumByUUID)
		albums.PUT("/:uuid", controllers.EditAlbum)
		albums.POST("/submit", controllers.CreateAlbum)
		albums.DELETE("/:uuid", controllers.DeleteAlbum)
		albums.PATCH("/images/:uuid", controllers.DeleteImageFromAlbum)
	}
}
