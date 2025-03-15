package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func LoadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return err
	}

	return nil
}

func GetString(key ConfigKey, defaultValue string) string {
	val := os.Getenv(string(key))
	if val == "" {
		return defaultValue
	}

	return val
}

func GetInt(key ConfigKey, defaultValue int) int {
	val := os.Getenv(string(key))
	if val == "" {
		return defaultValue
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		log.Println("error when parse config key", key, "with error", err.Error())
		return defaultValue
	}

	return valInt
}
