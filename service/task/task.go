package task

import (
	"axe-backend/store"
	"axe-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func QueryTasks(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	task, err := store.GetTaskByID(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"task": task})
}

func UpdateTask(c *gin.Context) {
	var req store.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.UpdateTs = util.GetCurrentTimestamp()
	err := req.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := store.DeleteTaskByID(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "task deleted successfully"})
}

func CreateTask(c *gin.Context) {
	var req store.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.CreateTs = util.GetCurrentTimestamp()
	req.UpdateTs = util.GetCurrentTimestamp()
	err := req.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"task": req, "tid": req.ID}).Info("task created successfully")
	c.JSON(200, gin.H{"message": "task created successfully", "task": req})
}