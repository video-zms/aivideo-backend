package route

import (
	projectService "axe-backend/service/project"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r gin.IRouter) {
	projectApi := r.Group("project")
	projectApi.POST("query", projectService.QueryProjects)
	projectApi.POST("update", projectService.UpdateProject)
	projectApi.POST("delete", projectService.DeleteProject)
	projectApi.POST("create", projectService.CreateProject)
}
