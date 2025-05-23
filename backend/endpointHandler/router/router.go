package router

import (
	analyticsHandler "github.com/endpointhandler/analytics"
	compareHandler "github.com/endpointhandler/compare"
	"github.com/endpointhandler/config"
	"github.com/endpointhandler/handler"
	"time"

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
			analytics.GET("/time-open", analyticsHandler.TimeOpenAnalytics)
			analytics.GET("/status-distribution", analyticsHandler.StatusDistribution)
			analytics.GET("/time-spent", analyticsHandler.TimeSpentAnalytics)
			analytics.GET("/priority", analyticsHandler.PriorityAnalytics)
		}

		compare := api.Group("/compare")
		{
			compare.GET("/time-open", compareHandler.CompareTimeOpen)
			compare.GET("/status-distribution", compareHandler.CompareStatusDistribution)
			compare.GET("/time-spent", compareHandler.CompareTimeSpent)
			compare.GET("/priority", compareHandler.ComparePriority)
		}
	}

	return r
}
