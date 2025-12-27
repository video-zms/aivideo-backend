package asset

import (
	assetService "axe-backend/service/asset"
	"github.com/gin-gonic/gin"
)


func SetupRouter(r gin.IRouter) {
	assetApi := r.Group("asset")
	assetApi.POST("upload", assetService.UploadAsset)
	assetApi.POST("query", assetService.QueryAssets)
	assetApi.POST("delete", assetService.DeleteAsset)
	assetApi.POST("update", assetService.UpdateAsset)
}