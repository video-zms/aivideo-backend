package asset

import (
	"axe-backend/store"
	"axe-backend/util"

	"github.com/gin-gonic/gin"
)

func UploadAsset(c *gin.Context) {
	var req store.Asset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := req.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "asset uploaded successfully", "asset": req})
}

func QueryAssets(c *gin.Context) {
	var req struct {
		ID   int64 `json:"id"`
		Type int64 `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.ID != 0 {
		asset, err := store.GetAssetInfo(req.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"assets": []interface{}{asset}})
		return
	}
	if req.Type != 0 {
		assets, err := store.GetAssetsByType(req.Type)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"assets": assets})
		return
	}
	c.JSON(400, gin.H{"error": "invalid request"})
}

func DeleteAsset(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	assetInfo, err := store.GetAssetInfo(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = assetInfo.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "asset deleted successfully"})
}

func UpdateAsset(c *gin.Context) {
	var req store.Asset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	assetInfo, err := store.GetAssetInfo(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if req.AssetType != 0 {
		assetInfo.AssetType = req.AssetType
	}
	if req.Detail != "" {
		assetInfo.Detail = req.Detail
	}
	if req.Extra != "" {
		assetInfo.Extra = req.Extra
	}
	assetInfo.UpdateTs = util.GetCurrentTimestamp()
	err = assetInfo.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "asset updated successfully"})
}
