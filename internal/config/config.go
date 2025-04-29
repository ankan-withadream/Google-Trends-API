package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	API_PORT     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var APP_CONFIG = Config{
	API_PORT:     8080,
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

const GOOGLE_TRENDS_BASE_URL = "https://trends.google.com/trending"

func Init_env() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error getting env files: ", err)
	}
}
