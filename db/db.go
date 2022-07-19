package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func Connect() *sql.DB {
	USER := os.Getenv("DBUSER")
	PASS := os.Getenv("DBPASS")
	HOST := os.Getenv("DBHOST")
	PORT := 3306
	DBNAME := os.Getenv("DBNAME")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		USER,
		PASS,
		HOST,
		PORT,
		DBNAME,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err)
	}

	if err != nil {
		panic(err.Error())
	}

	return db
}
