package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func startApiServer() error {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ApiPort)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r.Run(addr)
}
