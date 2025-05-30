package routes

import (
	"github.com/charis16/luminor-golang-be/src/controllers"
	"github.com/charis16/luminor-golang-be/src/middleware"
	"github.com/gin-gonic/gin"
)

func FaqRoutes(rg *gin.RouterGroup) {
	faq := rg.Group("/faqs")
	faq.GET("/", controllers.GetPublishedFaqs)
	faq.Use(middleware.AdminRequireAuth(), middleware.RequireRole("admin"))
	{
		faq.GET("/lists", controllers.GetFaqs)
		faq.GET("/:uuid", controllers.GetFaqByUUID)
		faq.PUT("/:uuid", controllers.EditFaq)
		faq.POST("/submit", controllers.CreateFaq)
		faq.DELETE("/:uuid", controllers.DeleteFaq)
	}
}
