package storages

import (
	"LamodaTest/internal/registry"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	AddRoute       = "/storages/add"
	DeleteRoute    = "/storages/delete"
	AvailableRoute = "/storages/available"
	AllRoute       = "/storages/all"
	AccessStatus   = "/storages/access"
)

type Handler struct {
	registry *registry.Database
	log      logrus.FieldLogger
}

type Add struct {
}

func NewHandler(registry *registry.Database, log logrus.FieldLogger) *Handler {
	return &Handler{registry: registry, log: log}
}

func (h *Handler) Add(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		Available *bool  `json:"available" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("can't parse body from `/storage/add` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	addedId, err := h.registry.StoragesAdd(context.Background(), input.Name, *input.Available)
	if err != nil {
		h.log.Errorf("can't add storage: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Not added"})
		return
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": addedId,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	var input struct {
		Id int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("can't parse body from `/storage/delete` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	deleted, err := h.registry.StoragesDelete(context.Background(), input.Id)
	if err != nil {
		h.log.Errorf("can't delete storage: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Can't delete this storage"})
		return
	}
	if deleted == 0 {
		c.JSON(200, gin.H{
			"code":    http.StatusOK,
			"message": "no records are deleted",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
	})
}

func (h *Handler) Available(c *gin.Context) {
	storages, err := h.registry.Storages(context.Background(), false)
	if err != nil {
		h.log.Errorf("can't get available storages: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": storages,
	})
}

func (h *Handler) All(c *gin.Context) {
	storages, err := h.registry.Storages(context.Background(), true)
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
	var input struct {
		Id        int   `json:"id" binding:"required"`
		Available *bool `json:"available" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("can't parse body from `/storage/add` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	changed, err := h.registry.StoragesChangeAccess(context.Background(), input.Id, *input.Available)
	if err != nil {
		h.log.Errorf("can't change storage: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Can't change this storage"})
		return
	}
	if changed == 0 {
		c.JSON(200, gin.H{
			"code":    http.StatusOK,
			"message": "no records are changed",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
	})
}
