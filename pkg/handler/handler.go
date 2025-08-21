package handler

import (
	_ "Test_project_Effective_Mobile/docs"
	"Test_project_Effective_Mobile/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subs := router.Group("/subs")
	{

		{
			subs.POST("/", h.createSubs)
			subs.GET("/", h.getAllSubs)
			subs.GET("/:id", h.getSubById)
			subs.PUT("/:id", h.updateSub)
			subs.DELETE("/:id", h.deleteSub)
			subs.GET("/total", h.getTotalPrice)

		}

	}
	return router

}
