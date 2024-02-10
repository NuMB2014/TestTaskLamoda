package handler

import (
	"LamodaTest/internal/handler/available"
	"LamodaTest/internal/handler/goods"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	reserveRoute      = "/reserve"
	releaseRoute      = "/release"
	getAvailableRoute = "/get_available"
)

func Router(log *logrus.Logger, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.Writer()))

	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.POST(reserveRoute, goods.Reserve)
	router.POST(releaseRoute, goods.Release)
	router.GET(getAvailableRoute, available.Get)
	return router
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "page not found"})
}

func notAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": http.StatusMethodNotAllowed, "message": "page not found"})
}
