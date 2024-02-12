package goods

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	//handler.Route
	log logrus.FieldLogger
}

func Add(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Delete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Release(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Reserve(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Remains(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func All(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}