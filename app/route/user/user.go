package user

import (
	"github.com/gin-gonic/gin"
	userService "axe-backend/service/user"
)

func SetupRouter(r gin.IRouter) {
	userApi := r.Group("user")
	userApi.POST("login", userService.Login)
	userApi.POST("register", userService.Register)
}
