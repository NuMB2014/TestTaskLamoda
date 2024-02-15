package main

import (
	"LamodaTest/internal/handler"
	"LamodaTest/internal/logger"
	"database/sql"
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

	db, err := sql.Open("mysql", getMysqlDSN())
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

	router := handler.Router(log, debug, db)
	err = router.Run(fmt.Sprintf("%s:%s", *ip, *port))
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

func getMysqlDSN() string { //"username[:password]@][protocol[(address)]]/dbname" root:1@lamoda_mysql/Lamoda
	database := os.Getenv("MYSQL_DATABASE")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, database)
	if len(str) <= 9 {
		return ""
	}
	return str
}
