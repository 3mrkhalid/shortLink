package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func GetMe(c *gin.Context) {

	user, err := c.Get("user")

	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

