package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

type webServiceConf struct {
	Host string
	Port string
}

func controllerIPLookUp(c *gin.Context) {
	nginxProxyIP := c.GetHeader("X-Real-IP")
	clientIP := c.ClientIP()

	ipNumber := clientIP
	if nginxProxyIP != "" {
		ipNumber = nginxProxyIP
	}

	c.JSON(http.StatusOK, gin.H{
		"ip": ipNumber,
	})
}

func controllerHTTPReqProxxy(c *gin.Context) {
	proxyReqURL, _ := c.GetQuery("url")
	proxyReqMethod, _ := c.GetQuery("method")

	if proxyReqMethod == "" {
		proxyReqMethod = "GET"
	}

	if proxyReqURL != "" {
		data := httpGetReq(proxyReqMethod, proxyReqURL, "http://code-http.jonham.cn")
		c.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"data": data,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "query.url is required",
		})
	}
}

func httpGetReq(method string, url string, proxyHost string) map[string]interface{} {
	timeout5s, _ := time.ParseDuration("5s")
	client := &http.Client{
		Timeout: timeout5s,
	}

	var req *http.Request
	req, _ = http.NewRequest(method, url, nil)

	req.Header.Add("Host", proxyHost)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result := make(map[string]interface{})

	result["Header"] = resp.Header
	body, err := ioutil.ReadAll(resp.Body)
	result["Body"] = string(body)

	return result
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// 替换默认的模板{{}}，避免与Vue默认语法冲突
	r.Delims("{[{", "}]}")

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ip-look-up", controllerIPLookUp)
		apiV1.POST("/req-proxy", controllerHTTPReqProxxy)
		apiV1.GET("/req-proxy", controllerHTTPReqProxxy)
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

func main() {
	data, _ := ioutil.ReadFile("config.yml")
	config := webServiceConf{}
	yaml.Unmarshal(data, &config)

	r := setupRouter()
	r.Run(config.Host + ":" + config.Port)
}
