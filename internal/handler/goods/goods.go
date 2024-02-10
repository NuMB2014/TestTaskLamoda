package goods

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	//handler.Route
	log logrus.FieldLogger
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
