package shot

import (
	"axe-backend/store"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func QueryShots(c *gin.Context) {
	var req struct {
		ShotID	string `json:"shot_id"`
		ID 	int64  `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.ID != 0 {
		shot, err := store.GetShotByID(req.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"shot": shot})
		return
	}
	if req.ShotID != "" {
		shot, err := store.GetShotByShotID(req.ShotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"shot": shot})
		return
	}
}

func CreateShot(c *gin.Context) {
	var req store.Shot
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := req.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"shot": req, "sid": req.ID}).Info("shot created successfully")
	c.JSON(200, gin.H{"message": "shot created successfully", "shot": req})
}

func UpdateShot(c *gin.Context) {
	var shot store.Shot
	if err := c.ShouldBindJSON(&shot); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := shot.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"shot": shot, "sid": shot.ID}).Info("shot updated successfully")
	c.JSON(200, gin.H{"message": "shot updated successfully"})
}

func DeleteShot(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	shot, err := store.GetShotByID(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if shot == nil {
		c.JSON(404, gin.H{"error": "shot not found"})
		return
	}
	err = shot.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"shot": shot, "sid": shot.ID}).Info("shot deleted successfully")
	c.JSON(200, gin.H{"message": "shot deleted successfully"})
}