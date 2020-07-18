package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setEnv(a *AmqpClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("amqpConnection", a)
		c.Next()
	}
}

func getAmqpConnection(c *gin.Context) (a *AmqpClient) {
	a, ok := c.Keys["amqpConnection"].(*AmqpClient)
	if !ok {
		log.Fatal("Connection to AMQP broker not set up")
	}
	return
}

func registerNewDeviceHandler(c *gin.Context) {
	aConnection := getAmqpConnection(c)
	var newRegister newRegister
	if err := c.ShouldBindJSON(&newRegister); err != nil {
		log.Println("Failed to bind json")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error"})
		return
	}
	res, code := aConnection.NewDeviceRPC(newRegister)
	log.Println(code)
	if code == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Couldn't get AMQP RPC status code"})
	} else if code == 400 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": res})
	} else if code == 200 {
		c.JSON(http.StatusOK, gin.H{"msg": res})
	}
}

func newReading(c *gin.Context) {
	log.Printf("%s", "TEST")
	c.Status(http.StatusOK)
}

func main() {
	r := gin.Default()
	a := AmqpClient{}
	a.ConnectToBroker("amqp://guest:guest@localhost:5672")
	r.Use(setEnv(&a))
	// Routes to use
	r.POST("/register_new_device", registerNewDeviceHandler)
	r.POST("/new_reading", newReading)

	r.Run("0.0.0.0:6000")
}
