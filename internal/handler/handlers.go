package handler

import (
	"LamodaTest/internal/handler/goods"
	"LamodaTest/internal/handler/storages"
	"LamodaTest/internal/registry"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"net/http"
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
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	reg := registry.New(db)
	goodH := goods.NewHandler(reg, log)
	storageH := storages.NewHandler(reg, log)
	reg.ReleaseGood(context.TODO(), 1, 1)
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
