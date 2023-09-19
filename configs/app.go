package config

import (
    _ "github.com/lib/pq"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "fmt"
)

var (
    db *gorm.DB
)

func Connect() *gorm.DB {
    connectionString := "postgres://lnrhlkqc:dSGqE8wF8kewswsubDnjurYmBStVStp0@jelani.db.elephantsql.com/lnrhlkqc"
    var err error
    db, err = gorm.Open("postgres", connectionString)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected to ElephantSQL")
    return db
}

func GetDB() *gorm.DB {
    return db
}
