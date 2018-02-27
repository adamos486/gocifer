package main

import (
	"cobb.com/server/database"
	"cobb.com/server/events/controllers"
	"cobb.com/server/events/services"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	//dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//open connection to db
	db, err := sql.Open("postgres", os.Getenv("DB_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	//define clients
	client := database.NewClient(db)
	service := services.NewClient(client)
	//init APIs
	controllers.NewEventsApiClient(service)

	router := gin.Default()
	router.POST("/add", controllers.AddEvent)
	router.GET("/events", controllers.GetEvents)
	router.Run(":9000")
}
