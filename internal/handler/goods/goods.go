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

type GoodWithCount struct {
	UniqId int `json:"uniq_id" binding:"required"`
	Count  int `json:"count" binding:"required"`
}

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
	var inputArr []GoodWithCount
	if err := c.ShouldBindJSON(&inputArr); err != nil {
		h.log.Errorf("can't parse body from `/good/release` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	var result []goods.ReleasedDTO
	for _, obj := range inputArr {
		err := h.registry.ReleaseGood(context.Background(), obj.UniqId, obj.Count)
		tmp := goods.ReleasedDTO{}
		tmp.UniqId = obj.UniqId
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
	var inputArr []GoodWithCount
	if err := c.ShouldBindJSON(&inputArr); err != nil {
		h.log.Errorf("can't parse body from `/good/reserve` request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid JSON"})
		return
	}
	var result []goods.ReservedDTO
	for _, obj := range inputArr {
		reserved, err := h.registry.ReserveGood(context.Background(), obj.UniqId, obj.Count)
		if err != nil {
			h.log.Warn(err)
		}
		tmp := goods.ReservedDTO{
			UniqId:         obj.UniqId,
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
	storages, err := h.registry.Goods(context.Background())
	if err != nil {
		h.log.Errorf("can't get all goods: %s", err.Error())
		c.JSON(500, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
	}
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": storages,
	})
}
