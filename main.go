package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 定义数据库模型
type Person struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	IDCard string `json:"id_card" db:"id_card"`
	Age    int    `json:"age" db:"age"`
	Gender string `json:"gender" db:"gender"`
}

func main1() {
	// 连接数据库
	dsn := "root:Gaochong806214@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化 Gin
	r := gin.Default()
	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")              // 允许所有来源
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS") // 允许的请求方法
		c.Header("Access-Control-Allow-Headers", "Content-Type")  // 允许的请求头

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 查询所有数据
	r.Any("/userinfo/query", func(c *gin.Context) {
		var people []Person
		query := "SELECT * FROM people"
		if err := db.Select(&people, query); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, people)
	})

	// 启动服务
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
