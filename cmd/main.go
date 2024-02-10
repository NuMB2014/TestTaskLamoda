package main

import (
	"LamodaTest/internal/handler"
	"LamodaTest/internal/logger"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	debug := isDebug()
	log := logger.New(debug)

	ip := flag.String("ip", "0.0.0.0", "ip address for web server")
	port := flag.String("port", "8080", "port for web server")
	flag.Parse()

	router := handler.Router(log, debug)
	err := router.Run(fmt.Sprintf("%s:%s", *ip, *port))
	if err != nil {
		log.Fatal(err)
	}
}

func isDebug() bool {
	debug := os.Getenv("DEBUG")
	parseBool, err := strconv.ParseBool(debug)
	if err != nil || !parseBool {
		return false
	}
	return true
}
