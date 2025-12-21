package main

import (
	"axe-backend/aivideo"
	"axe-backend/config"
	db "axe-backend/store"
	"axe-backend/user"
	"axe-backend/util"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"path/filepath"

	"github.com/cloudflare/tableflip"
	"github.com/gin-gonic/gin" // 替换 lfshook 为 rotatelogs
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	isDev      bool
	isLocalDev bool
	port       = "6869"
)

func init() {
	flag.BoolVar(&isDev, "is_dev", false, "isDev")
	flag.BoolVar(&isLocalDev, "is_local_dev", false, "isLocalDev")
	flag.Parse()

	// 创建 logs 目录
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// 配置 logrus 使用 rotatelogs 按日期分割文件
	logFile := filepath.Join(logDir, "app.%Y-%m-%d.log")
	rotateLogs, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 可选：保留7天
		rotatelogs.WithRotationTime(24*time.Hour), // 每天轮转
	)
	if err != nil {
		panic("Failed to create rotatelogs: " + err.Error())
	}
	logrus.SetOutput(rotateLogs)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	errfile, _ := os.OpenFile("./api_err.log", os.O_WRONLY|os.O_CREATE, 0755)
	syscall.Dup2(int(errfile.Fd()), 2)

	println("main start, is_dev:", isDev)

	//加载配置
	if isLocalDev {
		config.InitEnvConfCustom("./config/env.local.toml")
	} else if isDev {
		config.InitEnvConfCustom("../../config/env.dev.toml")
	} else {
		config.InitEnvConfCustom("./config/env.online.vpc.toml")
	}
}

func initPprof() {
	err := http.ListenAndServe(":22221", nil)
	if err != nil {
		println("start profile api err: ", err.Error())
		time.AfterFunc(time.Second*20, func() {
			err := http.ListenAndServe(":22221", nil)
			println("restart profile api err: ", err.Error())
		})
	}
}

func main() {
	// profile api
	go initPprof()
	var upg, _ = tableflip.New(tableflip.Options{PIDFile: "/tmp/api_server.pid"})
	defer upg.Stop()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2)
		for range sig {
			upg.Upgrade()
		}
	}()

	// redis连接
	// db.ConnectRedis()
	// 数据库连接
	db.ConnectMainDB()

	// SetMode
	if !config.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 路由配置
	r := setupRouter()

	println("main start, 监听端口:", util.GetLocalIp()+":"+port, ",PID:", strconv.Itoa(os.Getpid()))
	logrus.Infoln("main start, 监听端口:", util.GetLocalIp()+":"+port, ",PID:", strconv.Itoa(os.Getpid()))
	ln, err := upg.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	go server.Serve(ln)

	if err := upg.Ready(); err != nil {
		println("服务启动失败:" + err.Error())
		panic(err)
	}
	println("main start, 服务启动成功，监听端口:", port, ",PID:", strconv.Itoa(os.Getpid()))
	logrus.Infoln("main start, 服务启动成功，监听端口:", port, ",PID:", strconv.Itoa(os.Getpid()))
	<-upg.Exit()

	println("main start, 结束父进程 PID:", strconv.Itoa(os.Getpid()))
	logrus.Infoln("main start, 结束父进程 PID:", strconv.Itoa(os.Getpid()))
	time.AfterFunc(10*time.Second, func() {
		os.Exit(1) //防止 server.Shutdown 超时
	})
	server.Shutdown(context.Background())
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Any("checkup", func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().String())
	})

	aivideo.SetupRouter(r.Group("/aivideo"))
	user.SetupRouter(r.Group("/user"))

	return r
}
