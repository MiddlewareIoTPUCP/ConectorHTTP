package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/houseofcat/turbocookedrabbit/v2/pkg/tcr"
)

func setEnv(cp *tcr.ConnectionPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("amqpConnectionPool", cp)
		c.Next()
	}
}

func getConnectionPool(c *gin.Context) *tcr.ConnectionPool {
	cp, ok := c.Keys["amqpConnectionPool"].(*tcr.ConnectionPool)
	if !ok {
		log.Fatal("Connection Pool not set up")
	}
	return cp
}

func registerNewDeviceHandler(c *gin.Context) {
	cp := getConnectionPool(c)
	var newRegister newRegisterJSON
	if err := c.ShouldBindJSON(&newRegister); err != nil {
		log.Println("Failed to bind json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "JSON is not valid"})
		return
	}
	res, code := NewDeviceRPC(cp, newRegister)

	// We return depending on code
	switch code {
	case 0:
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Couldn't get AMQP RPC status code"})
	case 200:
		c.JSON(http.StatusOK, gin.H{"deviceID": res})
	case 400: // For validation errors
		c.JSON(http.StatusBadRequest, gin.H{"msg": res})
	case 500: // Couldn't connect to mongoDB
		c.JSON(http.StatusInternalServerError, gin.H{"msg": res})
	}
}

func newReadingHandler(c *gin.Context) {
	cp := getConnectionPool(c)
	var readings readingsJSON
	if err := c.ShouldBindJSON(&readings); err != nil {
		log.Println("Failed to bind JSON")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "JSON is not valid"})
		return
	}

	err := NewReading(cp, readings)
	if err != nil {
		log.Println("Couldn't send message to AMQP Broker")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "Couldn't send message to AMQP Broker"})
		return
	}
	c.Status(http.StatusOK)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	cp := ConnectToBroker("amqp://guest:guest@localhost:5672")
	r.Use(setEnv(cp))
	// Routes to use
	r.POST("/register_new_device", registerNewDeviceHandler)
	r.POST("/new_reading", newReadingHandler)
	r.Run("0.0.0.0:6000")
}
