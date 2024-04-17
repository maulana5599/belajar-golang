package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectionDatabase() {
	var err error
	DB, err = sql.Open("mysql", "maulana:maulana186@tcp(127.0.0.1:3306)/belajar")
	if err != nil {
		log.Fatalln(err.Error())
		panic(err.Error())
	}
	fmt.Println("Connecting to database")
}
