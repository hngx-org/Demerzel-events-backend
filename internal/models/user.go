package models

import (
	"demerzel-events/configs"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Person struct{
	gorm.Model
	Name string `gorm:""json:"name"`
}

func init() {
	config.Connect()
	db = config .GetDB()
	db.AutoMigrate(&Person{})
}

func GetDB() *gorm.DB {
    return db
}

func (b *Person) CreatePerson() *Person{
	db.NewRecord(b)
	db.Create(&b)
	return b
}


