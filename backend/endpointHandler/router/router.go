package router

import (
	"github.com/endpointhandler/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
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

		api.GET("/compare/:taskNumber", handler.CompareTask)
		api.GET("/compare/projects", handler.CompareAllProjects)
	}
	return r
}
