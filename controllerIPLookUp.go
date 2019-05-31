package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
