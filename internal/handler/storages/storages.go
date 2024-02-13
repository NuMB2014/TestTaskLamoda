package storages

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

func (h *Handler) Available(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) All(c *gin.Context) {
	storages, err := h.registry.Storages(context.Background())
	if err != nil {
		h.log.Errorf("can't get all storages: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": storages,
	})
}

func (h *Handler) ChangeAccess(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": "storages",
	})
}
