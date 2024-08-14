package connection

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"os"
)

func CreateDatabaseConnection() (*gorm.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")

	postgresConfig := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_PORT)

	fmt.Printf("postgres config => %s connection.go ", postgresConfig)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  postgresConfig,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	return db, err
}
