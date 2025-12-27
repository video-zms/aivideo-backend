package chapter

import (
	"axe-backend/store" 
	"axe-backend/util"

	"github.com/gin-gonic/gin"
)

func QueryChapters(c *gin.Context) {
	var req struct {
		ProjectID int64 `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chapters, err := store.ListChaptersByProjectID(req.ProjectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"chapters": chapters})
}

func UpdateChapter(c *gin.Context) {
	var req struct {
		ID        int64  `json:"id" binding:"required"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		ProjectID int64  `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chapter, err := store.GetChapterByID(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if req.Title != "" {
		chapter.StoryTitle = req.Title
	}
	if req.Content != "" {
		chapter.Story = req.Content
	}
	if req.ProjectID != 0 {
		chapter.ProjectID = req.ProjectID
	}
	err = chapter.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "chapter updated successfully"})
}

func DeleteChapter(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chapter, err := store.GetChapterByID(req.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = chapter.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "chapter deleted successfully"})
}

func CreateChapter(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Story   string `json:"story" binding:"required"`
		ProjectID int64  `json:"project_id" binding:"required"`
		StoryTitle string `json:"story_title"`
		StoryScene string `json:"story_scene"`
		StoryShots string `json:"story_shots"`
		Extea     string `json:"extea"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chapter := &store.Chapter{
		StoryTitle: req.Title,
		Story:   req.Story,
		ProjectID: req.ProjectID,
		StoryScene: req.StoryScene,
		StoryShots: req.StoryShots,
		Extea:     req.Extea,
		CreateTs:  util.GetCurrentTimestamp(),
		UpdateTs:  util.GetCurrentTimestamp(),
	}
	err := chapter.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "chapter added successfully", "chapter_id": chapter.ID})
}
