package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/register_new_device", func(c *gin.Context) {
		res := NewDeviceRPC()
		c.JSON(http.StatusOK, gin.H{"msg": res})
	})

	r.POST("/new_reading", func(c *gin.Context) {
		log.Printf("%s", "TEST")
		c.Status(http.StatusOK)
	})

	r.Run("0.0.0.0:6000")
}
