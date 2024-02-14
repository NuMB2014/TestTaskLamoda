package handler

import (
	"LamodaTest/internal/handler/goods"
	"LamodaTest/internal/handler/storages"
	"LamodaTest/internal/registry"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Router(log *logrus.Logger, debug bool, db *sql.DB) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.Writer()))

	reg := registry.New(db)
	goodH := goods.NewHandler(reg, log)
	storageH := storages.NewHandler(reg, log)
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.PUT(goods.AddRoute, goodH.Add)
	router.DELETE(goods.DeleteRoute, goodH.Delete)
	router.POST(goods.ReserveRoute, goodH.Reserve)
	router.POST(goods.ReleaseRoute, goodH.Release)
	router.GET(goods.RemainsRoute, goodH.Remains)
	router.GET(goods.AllRoute, goodH.All)

	router.PUT(storages.AddRoute, storageH.Add)
	router.DELETE(storages.DeleteRoute, storageH.Delete)
	router.GET(storages.AvailableRoute, storageH.Available)
	router.GET(storages.AllRoute, storageH.All)
	router.POST(storages.AccessStatus, storageH.ChangeAccess)

	return router
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "page not found"})
}

func notAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": http.StatusMethodNotAllowed, "message": "page not found"})
}
