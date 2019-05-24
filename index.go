package main

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

type webServiceConf struct {
	Port string
}

func controllerIPLookUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ip": c.ClientIP(),
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// 替换默认的模板{{}}，避免与Vue默认语法冲突
	r.Delims("{[{", "}]}")

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ip-look-up", controllerIPLookUp)
	}

	r.StaticFile("/", "./public/index.html")
	// r.StaticFile("/static", "./public/static")

	// 定义默认路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "Not Found",
		})
	})

	return r
}

func main() {
	data, _ := ioutil.ReadFile("config.yml")
	config := webServiceConf{}
	yaml.Unmarshal(data, &config)

	r := setupRouter()
	r.Run(":" + config.Port)
}
