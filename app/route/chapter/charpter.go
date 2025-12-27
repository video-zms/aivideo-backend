package chapter

import (
	charpterService "axe-backend/service/chapter"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r gin.IRouter) {
	charpterApi := r.Group("chapter")
	charpterApi.POST("query", charpterService.QueryChapters)
	charpterApi.POST("update", charpterService.UpdateChapter)
	charpterApi.POST("delete", charpterService.DeleteChapter)
	charpterApi.POST("create", charpterService.CreateChapter)
}
