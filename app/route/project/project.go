package route

import (
	"github.com/gin-gonic/gin"
	projectService "axe-backend/service/project"
)

func SetupRouter(r gin.IRouter) {
	projectApi := r.Group("project")
	projectApi.POST("query", projectService.QueryProjects);
	projectApi.POST("update", projectService.UpdateProject);
	projectApi.POST("delete", projectService.DeleteProject);
}