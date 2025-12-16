package db

import (
	"database/sql"
	"fmt"

	// _ "github.com/lib/pq" // pour postgres
	_ "github.com/go-sql-driver/mysql"
	// pour installer le driver
	// >go get github.com/.... (selon le package de driver à installer)
)

const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "api"
	password = "narciso"
	dbname   = "api"
)

var Conn *sql.DB

func NewDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// pour mysql
	//	var msqlInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
	//	 			user, password,host, port, dbname)

	conn, err := sql.Open(driver, psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connected to database !")
	return conn
}
