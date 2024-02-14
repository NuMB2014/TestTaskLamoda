package goods

import (
	"LamodaTest/internal/entity/goods"
	"LamodaTest/internal/registry"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	AddRoute     = "/goods/add"
	DeleteRoute  = "/goods/delete"
	ReserveRoute = "/goods/reserve"
	ReleaseRoute = "/goods/release"
	RemainsRoute = "/goods/remains"
	AllRoute     = "/goods/all"
)

type goodWithCount struct {
	UniqCode int `json:"uniq_code" binding:"required"`
	Count    int `json:"count" binding:"required"`
}

type Handler struct {
	registry registry.Db
	log      logrus.FieldLogger
}

func NewHandler(registry registry.Db, log logrus.FieldLogger) *Handler {
	return &Handler{registry: registry, log: log}
}

func (h *Handler) Add(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Size     string `json:"size" binding:"required"`
		UniqCode int    `json:"uniq_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("can't parse body from `/good/add` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	goodId, err := h.registry.GoodAdd(context.Background(), input.Name, input.Size, input.UniqCode)
	if err != nil {
		h.log.Errorf("can't add good: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Not added"})
		return
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": goodId,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	var input struct {
		UniqCode int `json:"uniq_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("can't parse body from `/good/delete` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	deleted, err := h.registry.GoodDelete(context.Background(), input.UniqCode)
	if err != nil {
		h.log.Errorf("can't delete good: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Can't delete this good"})
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

func (h *Handler) Release(c *gin.Context) {
	var inputArr []goodWithCount
	if err := c.ShouldBindJSON(&inputArr); err != nil {
		h.log.Errorf("can't parse body from `/good/release` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	var result []goods.ReleasedDTO
	for _, obj := range inputArr {
		err := h.registry.ReleaseGood(context.Background(), obj.UniqCode, obj.Count)
		tmp := goods.ReleasedDTO{}
		tmp.UniqCode = obj.UniqCode
		if err != nil {
			h.log.Warn(err)
			tmp.AdditionalInfo = "can't release this good"
		} else {
			tmp.AdditionalInfo = "OK"
		}
		result = append(result, tmp)
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": result,
	})
}

func (h *Handler) Reserve(c *gin.Context) {
	var inputArr []goodWithCount
	if err := c.ShouldBindJSON(&inputArr); err != nil {
		h.log.Errorf("can't parse body from `/good/reserve` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	var result []goods.ReservedDTO
	for _, obj := range inputArr {
		reserved, err := h.registry.ReserveGood(context.Background(), obj.UniqCode, obj.Count)
		if err != nil {
			h.log.Warn(err)
		}
		tmp := goods.ReservedDTO{
			UniqCode:       obj.UniqCode,
			Storages:       []map[string]int{},
			AdditionalInfo: "",
		}
		if len(reserved) > 0 {
			for storage, count := range reserved {
				tmp.Storages = append(tmp.Storages, map[string]int{
					"storage":  storage,
					"reserved": count,
				})
			}
			result = append(result, tmp)
		} else {
			tmp.AdditionalInfo = "Can't reserve this good"
			result = append(result, tmp)
		}
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": result,
	})
}

func (h *Handler) Remains(c *gin.Context) {
	list, err := h.registry.AvailableGoods(context.Background())
	if err != nil {
		h.log.Errorf("can't get available goods: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": list,
	})
}

func (h *Handler) All(c *gin.Context) {
	list, err := h.registry.Goods(context.Background())
	if err != nil {
		h.log.Errorf("can't get all goods: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": list,
	})
}
