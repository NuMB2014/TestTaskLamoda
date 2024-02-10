package available

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	//handler.Route
	log logrus.FieldLogger
}

func Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})

}
