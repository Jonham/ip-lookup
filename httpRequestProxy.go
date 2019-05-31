package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
