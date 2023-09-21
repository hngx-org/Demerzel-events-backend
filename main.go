package main

import (
    "demerzel-events/api"
    "demerzel-events/configs"
    "demerzel-events/services"
    "demerzel-events/internal/handlers/event_handler"
    "demerzel-events/internal/handlers/user_handler"  
    "demerzel-events/internal/models/event"
    "demerzel-events/internal/models/user"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "fmt"
    "os"
    "strconv"
    "time"
)

var db *gorm.DB

func main() {
    r := gin.Default()
    var err error
    db, err = setupDatabase()
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    r.POST("/events", CreateEvent)
    r.GET("/events", ListEvents)
    r.GET("/events/:id", GetEventByID)
    r.PUT("/events/:id", UpdateEvent)
    r.DELETE("/events/:id", DeleteEvent)
    r.POST("/users", CreateUser)
    r.GET("/users/:id", GetUserByID)
    r.PUT("/users/:id", UpdateUser)
    r.DELETE("/users/:id", DeleteUser)    

    r.Run(":8080") // Replace with the desired port
}


func (e *Event) CreateEvent(db *gorm.DB) error {
    // Database logic here
    return nil
}

// Placeholder function for checking if the user is authenticated and authorized
func isUserAuthenticatedAndAuthorized(c *gin.Context) bool {
    // Authentication and authorization logic here
    // middleware here
    return true // Return true if authenticated and authorized; otherwise, return false
    
}

// Placeholder function for getting the current user based on the authentication
func getCurrentUser(c *gin.Context) User {
    // Authentication and authorization logic here
    
    // For demonstration purposes, to create a sample user
    user := User{
        ID:       1,
        Username: "sampleuser",
        Email:    "sample@example.com",
        Password: "password",
        //
    }
    return user
}

func setupDatabase() (*gorm.DB, error) {
    // Replace with database connection details
    dsn := "Tipsyx:#MarquessROOT28@tcp(127.0.0.1:3306)/eventlist?charset=utf8&parseTime=True&loc=Local"

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("Failed to connect to the database: %v", err)
        return nil, err
    }

    if err := db.AutoMigrate(&Event{}, &User{}); err != nil {
        log.Printf("Failed to auto-migrate models: %v", err)
        return nil, err
    }
    
    return db, nil
}