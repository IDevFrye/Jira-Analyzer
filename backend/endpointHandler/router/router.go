package router

import (
	"github.com/endpointhandler/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Добавляем CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Разрешённые домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	{
		api.GET("/projects", handler.GetProjects)
		api.GET("/projects/:id", handler.GetProjectStats)
		api.DELETE("/projects/:id", handler.DeleteProject)

		connector := api.Group("/connector")
		{
			connector.GET("/projects", handler.GetJiraProjects)
			connector.POST("/updateProject", handler.UpdateJiraProject)
		}

		analytics := api.Group("/analytics")
		{
			analytics.GET("/time-open", handler.TimeOpenAnalytics)
			analytics.GET("/status-distribution", handler.StatusDistribution)
			analytics.GET("/time-spent", handler.TimeSpentAnalytics)
			analytics.GET("/priority", handler.PriorityAnalytics)
		}

		compare := api.Group("/compare")
		{
			compare.GET("/time-open", handler.CompareTimeOpen)
			compare.GET("/status-distribution", handler.CompareStatusDistribution)
			compare.GET("/time-spent", handler.CompareTimeSpent)
			compare.GET("/priority", handler.ComparePriority)
		}
	}

	return r
}
