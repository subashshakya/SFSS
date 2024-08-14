package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subashshakya/SFSS/db/connection"
	"github.com/subashshakya/SFSS/db/orms"
)

func main() {
	IP := os.Getenv("IP")
	PORT := os.Getenv("PORT")
	serverConfig := fmt.Sprintf("%s/%s", IP, PORT)
	fmt.Print(IP, PORT)
	db, err := connection.CreateDatabaseConnection()
	r := gin.Default()
	r.Use()
	if err == nil {
		fmt.Println("DB Connection Successful")
	} else {
		fmt.Println("Error => ", err)
	}
	orms.DatabaseConnection = db
	dbConn, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()
	serverRunErr := r.Run(serverConfig)
	if serverRunErr != nil {
		panic(serverRunErr)
	}
}
