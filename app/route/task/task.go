package task

import (
	taskService "axe-backend/service/task"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r gin.IRouter) {
	taskApi := r.Group("task")
	taskApi.POST("query", taskService.QueryTasks)
	taskApi.POST("update", taskService.UpdateTask)
	taskApi.POST("delete", taskService.DeleteTask)
	taskApi.POST("create", taskService.CreateTask)
}