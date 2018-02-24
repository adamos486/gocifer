package main

import (
	"cobb.com/server/database"
	"cobb.com/server/events/controllers"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"cobb.com/server/events/services"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	//dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_ADDR"))
	if err != nil {
		log.Fatal(err)
	}

	//define clients
	client := database.NewClient(db)
	service := services.NewClient(client)
	//init APIs
	controllers.NewEventsApiClient(client, service)

	router := gin.Default()
	router.POST("/add", controllers.AddEvent)
	router.GET("/events", controllers.GetEvents)
	router.Run(":9000")
}