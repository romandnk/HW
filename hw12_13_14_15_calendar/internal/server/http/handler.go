package internalhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
)

type Handler struct {
	*gin.Engine
	Services service.Services
	Logger   logger.Logger
}

func NewHandler(services service.Services, logger logger.Logger) *Handler {
	return &Handler{
		Services: services,
		Logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(LoggerMiddleware(h.Logger))
	gin.SetMode(gin.ReleaseMode)
	h.Engine = router

	api := router.Group("/api")
	{
		version := api.Group("/v1")
		{
			adverts := version.Group("/events")
			{
				adverts.POST("", h.CreateEvent)
				adverts.PATCH("/:id", h.UpdateEvent)
				adverts.DELETE("/:id", h.DeleteEvent)
				//adverts.GET("/:date", h.GetAllByDayEvents)
				//adverts.GET("/:date", h.GetAllByWeekEvents)
				//adverts.GET("/:date", h.GetAllByMonthEvents)
			}
		}
	}

	return router
}
