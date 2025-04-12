package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sakkyoi/DDDNS/internal/middleware"
	"net/http"
)

func startApiServer() error {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ApiPort)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r.Run(addr)
}
