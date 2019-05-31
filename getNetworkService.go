package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ipLogger struct {
	IP         string
	UpdateTime string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// 获取Get网络公网IP，并提交
func getNetIP() map[string]interface{} {
	ipLookupServiceURL := "https://ip-lookup.jonham.app/api/v1/save-get-office-ip"
	// ipLookupServiceURL := "http://localhost:19524/api/v1/save-get-office-ip"

	timeout5s, _ := time.ParseDuration("5s")
	client := &http.Client{
		Timeout: timeout5s,
	}
	var req *http.Request
	req, _ = http.NewRequest("GET", ipLookupServiceURL, nil)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	check(err)

	return result
}

func controllerGetNetIP(c *gin.Context) {
	JSONParse := newJSONStruct()
	d := ipLogger{}
	JSONParse.Load("./get-network-ip.log.json", &d)

	fmt.Println(d)

	c.JSON(http.StatusOK, gin.H{
		"ip":         d.IP,
		"updateTime": d.UpdateTime,
	})
}

func controllerSaveGetNetIP(c *gin.Context) {
	nginxProxyIP := c.GetHeader("X-Real-IP")
	clientIP := c.ClientIP()

	ipNumber := clientIP
	if nginxProxyIP != "" {
		ipNumber = nginxProxyIP
	}

	now, _ := time.Now().MarshalText()
	data := ipLogger{
		IP:         ipNumber,
		UpdateTime: string(now),
	}
	jsonStr, _ := json.Marshal(data)

	err := ioutil.WriteFile("./get-network-ip.log.json", jsonStr, 0755)
	check(err)

	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
		"ip":  ipNumber,
	})
}

func sendGetNetIP() {
	result := getNetIP()
	fmt.Println(result, " on ", time.Now())
}

// intervalMin 多少分钟执行一次
func startInterval(f func(), intervalMin time.Duration, intervalSecond time.Duration) {
	go func() {
		for {
			f()
			now := time.Now()
			next := now.Add(time.Minute*intervalMin + time.Second*intervalSecond)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}
