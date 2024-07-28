package main

import (
	"fmt"

	"github.com/subashshakya/SFSS/db/connection"
)

func main() {
	_, err := connection.CreateDatabaseConnection()
	if err != nil {
		fmt.Print("DB Connection Successful")
	}
}
