package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DataBaseConfig struct {
	Url string
}

type LogConfig struct {
	Level  int
	Format string
}

type MQTTConfig struct {
	Port               int
	Broker             string
	StatusSendInterval time.Duration
	QOSLevel           int
}

type RedisConfig struct {
	Port        int           `yaml:"port"`
	Url         string        `yaml:"url"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}
	log.Println(".env file loaded")
}

func NewMQTTConfig() *MQTTConfig {
	return &MQTTConfig{
		Port:               getInt("MQTT_PORT", 1883),
		Broker:             getString("MQTT_URL", "tcp://127.0.0.1:1883"),
		StatusSendInterval: getTimeDuration("STATUS_SEND_INTERVAL", 30),
		QOSLevel:           getInt("QOS_LEVEL", 1),
	}
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Port:        getInt("REDIS_PORT", 6379),
		Url:         getString("REDIS_URL", "127.0.0.1"),
		Password:    getString("REDIS_PASS", "my_pass"),
		User:        getString("REDIS_USER", "user"),
		DB:          getInt("REDIS_DATABASE", 0),
		MaxRetries:  getInt("REDIS_MAXRET", 5),
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
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

func getTimeDuration(key string, defValue int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(defValue) * time.Second
	}
	duration, err := time.ParseDuration(value)
	if err == nil {
		return duration
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(defValue) * time.Second
	}
	return time.Duration(i) * time.Second
}
