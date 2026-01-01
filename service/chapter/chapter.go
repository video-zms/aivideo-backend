package chapter

import (
	"axe-backend/store"
	"axe-backend/util"

	"github.com/gin-gonic/gin"
)

func QueryChapters(c *gin.Context) {
	var req struct {
		ID        int64 `json:"id"`
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
		ID         int64  `json:"id" binding:"required"`
		StoryTitle string `json:"story_title" `
		Story      string `json:"story"`
		ProjectID  int64  `json:"project_id" `
		StoryScene string `json:"story_scene"`
		StoryShots string `json:"story_shots"`
		Extra      string `json:"extra"`
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
	if req.StoryTitle != "" {
		chapter.StoryTitle = req.StoryTitle
	}
	if req.Story != "" {
		chapter.Story = req.Story
	}
	if req.ProjectID != 0 {
		chapter.ProjectID = req.ProjectID
	}
	if req.StoryScene != "" {
		chapter.StoryScene = req.StoryScene
	}
	if req.StoryShots != "" {
		chapter.StoryShots = req.StoryShots
	}
	if req.Extra != "" {
		chapter.Extra = req.Extra
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
	if chapter == nil {
		c.JSON(404, gin.H{"error": "chapter not found"})
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
		StoryTitle string `json:"story_title" binding:"required"`
		Story      string `json:"story" binding:"required"`
		ProjectID  int64  `json:"project_id" binding:"required"`
		StoryScene string `json:"story_scene"`
		StoryShots string `json:"story_shots"`
		Extra      string `json:"extra"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chapter := &store.Chapter{
		StoryTitle: req.StoryTitle,
		Story:      req.Story,
		ProjectID:  req.ProjectID,
		StoryScene: req.StoryScene,
		StoryShots: req.StoryShots,
		Extra:      req.Extra,
		CreateTs:   util.GetCurrentTimestamp(),
		UpdateTs:   util.GetCurrentTimestamp(),
	}
	err := chapter.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "chapter added successfully", "chapter_id": chapter.ID})
}
