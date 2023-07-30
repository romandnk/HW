package internalhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
)

type HandlerHTTP struct {
	engine   *gin.Engine
	services service.Services
	logger   logger.Logger
}

func NewHandlerHTTP(services service.Services, logger logger.Logger) *HandlerHTTP {
	return &HandlerHTTP{
		services: services,
		logger:   logger,
	}
}

func (h *HandlerHTTP) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(loggerMiddleware(h.logger))
	gin.SetMode(gin.ReleaseMode)
	h.engine = router

	api := router.Group("/api")
	{
		version := api.Group("/v1")
		{
			adverts := version.Group("/events")
			{
				adverts.POST("", h.CreateEvent)
				adverts.PATCH("/:id", h.UpdateEvent)
				adverts.DELETE("/:id", h.DeleteEvent)
				adverts.GET("/day/:date", h.GetAllByDayEvents)
				adverts.GET("/week/:date", h.GetAllByWeekEvents)
				adverts.GET("/month/:date", h.GetAllByMonthEvents)
			}
		}
	}

	return router
}
