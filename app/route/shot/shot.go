package shot


import (	"github.com/gin-gonic/gin"
	shotService "axe-backend/service/shot"
)

func SetupRouter(r gin.IRouter) {
	shotApi := r.Group("shot")
	shotApi.POST("create", shotService.CreateShot)
	shotApi.POST("query", shotService.QueryShots)
	shotApi.POST("update", shotService.UpdateShot)
	shotApi.POST("delete", shotService.DeleteShot)
}