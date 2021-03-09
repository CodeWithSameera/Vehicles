package helpers

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func ConnectDB() *gorm.DB {
	username := GoDotEnvVariable("APP_DB_USERNAME")
	password := GoDotEnvVariable("APP_DB_PASSWORD")
	dbname := GoDotEnvVariable("APP_DB_NAME")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user="+username+" password="+password+" dbname="+dbname+" port=5432 sslmode=disable TimeZone=Asia/Colombo",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}

	return db
}
