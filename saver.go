package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	driver       string
	dbConnection string
	dbUser       string
	dbPass       string
	endpoint     string
	instance     string
	options      string
)

func InitFlag() {
	flag.StringVar(&driver, "driver", "mysql", "specify database driver will be used")
	flag.StringVar(&dbUser, "dbUser", "root", "")
	flag.StringVar(&dbPass, "dbPass", "Admin#1234", "")
	flag.StringVar(&endpoint, "endpoint", "192.168.123.24:3306", "database ip:port")
	flag.StringVar(&instance, "instance", "mini4k", "specify database instance will be used")
	flag.StringVar(&options, "options", "", "database connection options")
}

func NewDatabaseConnection() (*sql.DB, error) {
	switch driver {
	case "mysql":
		dbConnection = fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", dbUser, dbPass, endpoint, instance, options)
	default:
		dbConnection = ""
	}
	db, err := sql.Open(driver, dbConnection)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db, nil
}
