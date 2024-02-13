package goods

import (
	"LamodaTest/internal/registry"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	registry *registry.Database
	log      logrus.FieldLogger
}

func NewHandler(registry *registry.Database, log logrus.FieldLogger) *Handler {
	return &Handler{registry: registry, log: log}
}

func (h *Handler) Add(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) Delete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) Release(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) Reserve(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) Remains(c *gin.Context) {
	goods, err := h.registry.AvailableGoods(context.Background())
	if err != nil {
		h.log.Errorf("can't get available goods: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": goods,
	})
}

func (h *Handler) All(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
