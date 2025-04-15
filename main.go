package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"jwt-banking/database"

	"jwt-banking/routes"
)

func main() {
	r := gin.Default()

	_, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to the database")
	}

	r.Use(cors.Default())
	r.Use(gin.Logger())

	routes.SetUpRoutes(r)
	fmt.Println("connected to the database")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server Failed to Run")
		return
	} else {
		fmt.Println("Server started successfully ")
	}

}
