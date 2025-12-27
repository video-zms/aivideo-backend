package project

import (
	projectStore "axe-backend/store"
	"axe-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func QueryProjects(c *gin.Context) {
	var req struct {
		Id int64 `json:"id"`
		Creator string `json:"creator"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Id == 0 {
		projects, err := projectStore.ListAllProjects()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"projects": projects})
		return
	}
	if req.Id != 0 {
		project,  err := projectStore.GetProjectById(req.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"project": project})
	}
	if req.Creator != "" {
		projects, err := projectStore.GetProjectsByCreator(req.Creator)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"projects": projects})
		return
	}
	c.JSON(400, gin.H{"error": "invalid request"})
}

func UpdateProject(c *gin.Context) {
	var req struct{
		Id int64 `json:"id" binding:"required"`
		Desc string `json:"desc"`
		Extea string `json:"extea"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	projectInfo, err := projectStore.GetProjectById(req.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if projectInfo == nil {
		c.JSON(404, gin.H{"error": "project not found"})
		return
	}
	if req.Desc != "" {
		projectInfo.Desc = req.Desc
	}
	if req.Extea != "" {
		projectInfo.Extea = req.Extea
	}
	projectInfo.UpdateTs = util.GetCurrentTimestamp()
	err = projectInfo.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"project": projectInfo, "pid": projectInfo.ID}).Info("project updated successfully")
	c.JSON(200, gin.H{"message": "project updated successfully"})
}


func DeleteProject(c *gin.Context) {
	var req struct{
		Id int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	projectInfo, err := projectStore.GetProjectById(req.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if projectInfo == nil {
		c.JSON(404, gin.H{"error": "project not found"})
		return
	}
	err = projectInfo.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{"project": projectInfo, "pid": projectInfo.ID}).Info("project deleted successfully")
	c.JSON(200, gin.H{"message": "project deleted successfully"})
}