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

const (
	goodsAddRoute     = "/goods/add"
	goodsDeleteRoute  = "/goods/delete"
	goodsReserveRoute = "/goods/reserve"
	goodsReleaseRoute = "/goods/release"
	goodsRemainsRoute = "/goods/remains"
	goodsAllRoute     = "/goods/all"

	storageAddRoute       = "/storages/add"
	storageDeleteRoute    = "/storages/delete"
	storageAvailableRoute = "/storages/available"
	storageAllRoute       = "/storages/all"
	storageAccessStatus   = "/storages/access"
)

func Router(log *logrus.Logger, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.Writer()))

	db, err := sql.Open("mysql", "root:1@/Lamoda")
	if err != nil {
		log.Fatalf("Can't connect to mysql: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't ping mysql: %v", err)
	}
	reg := registry.New(db)
	goodH := goods.NewHandler(reg, log)
	storageH := storages.NewHandler(reg, log)
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.GET(goodsAddRoute, goodH.Add)
	router.DELETE(goodsDeleteRoute, goodH.Delete)
	router.POST(goodsReserveRoute, goodH.Reserve)
	router.POST(goodsReleaseRoute, goodH.Release)
	router.GET(goodsRemainsRoute, goodH.Remains)
	router.GET(goodsAllRoute, goodH.All)

	router.GET(storageAddRoute, storageH.Add)
	router.DELETE(storageDeleteRoute, storageH.Delete)
	router.GET(storageAvailableRoute, storageH.Available)
	router.GET(storageAllRoute, storageH.All)
	router.POST(storageAccessStatus, storageH.ChangeAccess)

	return router
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "page not found"})
}

func notAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": http.StatusMethodNotAllowed, "message": "page not found"})
}
