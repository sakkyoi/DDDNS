package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sakkyoi/DDDNS/internal/middleware"
)

func startApiServer() error {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ApiPort)

	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// API Handlers
	api := r.Group("/api")
	api.POST("/register", register)
	api.DELETE("/register", unregister)

	// Fallback
	r.NoRoute(fallback)

	return r.Run(addr)
}
