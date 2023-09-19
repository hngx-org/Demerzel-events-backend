package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	var (
		dbUsername = os.Getenv("MYSQL_USERNAME")
		dbPass     = os.Getenv("MYSQL_PASSWORD")
		dbHost     = os.Getenv("MYSQL_HOST")
		dbName     = os.Getenv("MYSQL_DBNAME")
		dbPortStr  = os.Getenv("MYSQL_PORT")
	)

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert MYSQL_PORT to an integer: %v", err))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUsername, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to open DB: %v", err))
	}

	if gin.Mode() == gin.ReleaseMode {
		db.Logger.LogMode(0)
	}

	DB = db
	rawDB := RawDB()

	rawDB.SetMaxIdleConns(20)
	rawDB.SetMaxOpenConns(100)

	Migrate()
}

// RawDB returns the raw SQL database instance.
func RawDB() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}

	return db
}
