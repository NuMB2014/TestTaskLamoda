package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type db interface {
	Storages()
	Goods()
	ReserveGoods()
	ReleaseGoods()
}

type Database struct {
	db *sql.Conn
}

func New(connect *sql.Conn) *Database {
	return &Database{db: connect}
}

func (d *Database) Storages() {

}

func (d *Database) Goods() {

}

func (d *Database) ReserveGoods() {

}

func (d *Database) ReleaseGoods() {

}
