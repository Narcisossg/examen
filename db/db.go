package db

import (
	"database/sql"
	"fmt"

	// _ "github.com/lib/pq" // pour postgres
	// _ "github.com/go-sql-driver/mysql"
	// pour installer le driver
	// >go get github.com/.... (selon le package de driver à installer)
	//pour oracle
	_ "github.com/sijms/go-ora/v2"
)

const (
	driver      = "oracle"
	host        = "localhost"
	port        = 1521
	user        = "prepExamen"
	password    = "narciso"
	serviceName = "XE"
)

var Conn *sql.DB

func NewDB() *sql.DB {
	oracleInfo := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		user, password, host, port, serviceName)

	conn, err := sql.Open(driver, oracleInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connected to database !")
	return conn
}
