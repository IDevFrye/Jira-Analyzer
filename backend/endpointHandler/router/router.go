package router

import (
	"github.com/endpointhandler/config"
	"time"

	"github.com/endpointhandler/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://frontend:3000"}, // Разрешённые домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	{
		api.GET("/projects", func(c *gin.Context) {
			handler.GetProjects(c, cfg)
		})
		api.GET("/projects/:id", handler.GetProjectStats)
		api.DELETE("/projects/:id", handler.DeleteProject)

		connector := api.Group("/connector")
		{
			connector.GET("/projects", func(c *gin.Context) {
				handler.GetJiraProjects(c, cfg)
			})
			connector.POST("/updateProject", func(c *gin.Context) {
				handler.UpdateJiraProject(c, cfg)
			})
		}

		analytics := api.Group("/analytics")
		{
			analytics.GET("/time-open", handler.TimeOpenAnalytics)            // param = projectKey. если ишью открыта в диапазоне дней то количество ишью на одном проекте/ range count
			analytics.GET("/status-distribution", handler.StatusDistribution) // ля каждого статуса количество/ status count
			analytics.GET("/time-spent", handler.TimeSpentAnalytics)          //время затраченное на проект для каждого автора
			analytics.GET("/priority", handler.PriorityAnalytics)
		}

		compare := api.Group("/compare")
		{
			compare.GET("/time-open", handler.CompareTimeOpen) //key
			compare.GET("/status-distribution", handler.CompareStatusDistribution)
			compare.GET("/time-spent", handler.CompareTimeSpent)
			compare.GET("/priority", handler.ComparePriority)
		}
	}

	return r
}
