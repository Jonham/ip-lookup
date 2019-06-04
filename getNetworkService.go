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

func check(e error) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in check", r)
		}
	}()

	if e != nil {
		panic(e)
	}

	return "GOOD"
}

// 获取Get网络公网IP，并提交
func getNetIP() map[string]interface{} {
	// 处理错误信息
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getNetIP", r)
		}
	}()

	ipLookupServiceURL := "https://ip-lookup.jonham.app/api/v1/save-get-office-ip"
	// error
	// ipLookupServiceURL := "https://ip-lookupd.jonam.ee/api/error"
	// ipLookupServiceURL := "http://localhost:19524/api/v1/save-get-office-ip"

	timeout5s, _ := time.ParseDuration("5s")
	client := &http.Client{
		Timeout: timeout5s,
	}
	var req *http.Request
	req, err := http.NewRequest("GET", ipLookupServiceURL, nil)
	check(err)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("net ", err)
		check(err)
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
	// 处理错误信息
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in sendGetNetIP", r)
		}
	}()

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
