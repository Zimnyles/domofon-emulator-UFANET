package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DataBaseConfig struct {
	Url string
}

type LogConfig struct {
	Level  int
	Format string
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}
	log.Println(".env file loaded")
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}

func NewDBConfig() *DataBaseConfig {
	databaseURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),    
		os.Getenv("DB_USER"),    
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME"),     
		os.Getenv("DB_PORT"))
	return &DataBaseConfig{
		Url: databaseURL,
	}
}

func getString(key, defValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defValue
	}
	return value
}

func getInt(key string, defValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defValue
	}
	if value == "" {
		return defValue
	}
	return i
}
