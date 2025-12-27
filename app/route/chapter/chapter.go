package chapter

import (
	chapterService "axe-backend/service/chapter"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r gin.IRouter) {
	chapterApi := r.Group("chapter")
	chapterApi.POST("query", chapterService.QueryChapters)
	chapterApi.POST("update", chapterService.UpdateChapter)
	chapterApi.POST("delete", chapterService.DeleteChapter)
	chapterApi.POST("create", chapterService.CreateChapter)
}
