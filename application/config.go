package application

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort uint16
	DBUsername string
	DBPassword string
	DBName     string
	DBPort     string
	ESUserName string
	ESPassword string
	ESPort     string
}

func LoadConfig() Config {
	cfg := Config{}

	if serverPort := getDotEnvValue("SERVER_PORT"); serverPort != "" {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	cfg.DBName = getDotEnvValue("DB_NAME")
	cfg.DBUsername = getDotEnvValue("DB_USERNAME")
	cfg.DBPassword = getDotEnvValue("DB_PASSWORD")
	cfg.DBPort = getDotEnvValue("DB_PORT")
	cfg.ESPassword = getDotEnvValue("ELASTIC_PASSWORD")
	cfg.ESUserName = getDotEnvValue("ELASTIC_USERNAME")

	return cfg
}

// func getEnvValue(envName string, def string) string {
// 	if val, exists := os.LookupEnv(envName); exists {
// 		return val
// 	}
// 	return def
// }

func getDotEnvValue(key string) string {
	err := godotenv.Load("/go/src/app/.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return os.Getenv(key)
}
