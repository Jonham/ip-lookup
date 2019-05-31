package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// 替换默认的模板{{}}，避免与Vue默认语法冲突
	r.Delims("{[{", "}]}")

	apiV1 := r.Group("/api/v1")
	{
		// ip查找服务 https://ip-lookup.jonham.app
		apiV1.GET("/ip-look-up", controllerIPLookUp)

		// 服务于 https://code.jonham.app/cors/ 跨域的http请求代理
		apiV1.POST("/req-proxy", controllerHTTPReqProxxy)
		apiV1.GET("/req-proxy", controllerHTTPReqProxxy)

		// 服务于Get路由器IP查找
		apiV1.GET("/get-office-ip", controllerGetNetIP)
		apiV1.GET("/save-get-office-ip", controllerSaveGetNetIP)
	}

	r.StaticFile("/", "./public/index.html")

	// 定义默认路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "Not Found",
		})
	})

	return r
}
