package main

import (
	"slip/common"
	"slip/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	common.InitDB() //初始数据库

	r := gin.Default() //默认路由
	r = routers.CollectRoute(r)
	panic(r.Run()) // 监听并在 0.0.0.0:8080 上启动服务
}
