package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Domain string `json:"domain" binding:"required"`
	DestIp string `json:"dest_ip" binding:"required"`
}

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// register to the store
	var domain string
	if req.Domain == "." {
		domain = fmt.Sprintf("%s.", cfg.Domain)
	} else {
		domain = fmt.Sprintf("%s.%s.", req.Domain, cfg.Domain)
	}

	err := s.Register(domain, c.ClientIP(), req.DestIp, 0) // TODO: set ttl

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to register",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

type UnregisterRequest struct {
	Domain string `json:"domain" binding:"required"`
}

func unregister(c *gin.Context) {
	var req UnregisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := s.Unregister(req.Domain, c.ClientIP())
	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to unregister",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func fallback(c *gin.Context) {
	sourceIp := c.ClientIP()

	c.String(http.StatusOK, fmt.Sprintf(
		"This is a fallback page of DDDNS. \n\n"+
			"Why Your're Here? \n"+
			"1. Your DNS (Domain Name System) server didn't pass your information to DDDNS server due to it's privacy policy. (e.g. Cloudflare's 1.1.1.1) \n"+
			"2. There is no record registered for your IP. \n"+
			"3. The DNS server modified our response for some efficiency reasons. \n"+
			"4. Your device cached the old DNS response (e.g. your browser, OS, etc.) \n\n"+
			"Your IP: %s",
		sourceIp))
}
