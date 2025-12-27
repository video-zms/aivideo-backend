package user

import (
	userStore "axe-backend/store"
	"axe-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email	string `json:"email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 用户名或邮箱不能为空
	if req.Username == "" && req.Email == "" {
		c.JSON(400, gin.H{"error": "username or email is required"})
		return
	}
	user, err := userStore.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if user == nil {
		c.JSON(401, gin.H{"error": "user not found"})
		return
	}
	// 验证密码
	if user.Password != req.Password {
		c.JSON(401, gin.H{"error": "invalid password"})
		return
	}
	// 登录成功，返回用户信息（不包含密码）
	logrus.WithFields(logrus.Fields{"user": user, "uid": user.ID}).Info("user logged in successfully")
	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"create_ts": user.CreateTs,
		"update_ts": user.UpdateTs,
		"privilege": user.Privilege,
		"coin":      user.Coin,
	})
}

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Privilege int    `json:"privilege"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 检查用户名或邮箱是否已存在
	existingUser, err := userStore.GetUserByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if existingUser != nil {
		c.JSON(409, gin.H{"error": "username or email already exists"})
		return
	}
	// 创建新用户
	newUser := &userStore.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	newUser.CreateTs = util.GetCurrentTimestamp()
	newUser.UpdateTs = newUser.CreateTs
	newUser.Privilege = req.Privilege
	err = newUser.Add()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}
	logrus.WithFields(logrus.Fields{"user": newUser, "uid": newUser.ID}).Info("user registered successfully")
	c.JSON(201, gin.H{
		"id":       newUser.ID,
		"username": newUser.Username,
		"email":    newUser.Email,
		"create_ts": newUser.CreateTs,
		"update_ts": newUser.UpdateTs,
		"privilege": newUser.Privilege,
		"coin":      newUser.Coin,
	})
}