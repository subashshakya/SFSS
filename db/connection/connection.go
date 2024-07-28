package connection

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
)

func CreateDatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=172.17.0.2 user=postgres password=12345678 port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	fmt.Print(db)
	fmt.Printf("\n")
	return db, err
}
