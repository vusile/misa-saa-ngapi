package application

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort uint16
	DBUsername string
	DBPassword string
	DBName     string
	DBPort     string
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort: 3000,
		DBUsername: "root",
		DBPassword: "123456",
		DBName:     "masstimings",
		DBPort:     "db:3306",
	}

	if serverPort, exists := os.LookupEnv("MT_SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	cfg.DBName = getEnvValue("MT_DB_NAME", cfg.DBName)
	cfg.DBUsername = getEnvValue("MT_DB_USERNAME", cfg.DBUsername)
	cfg.DBPassword = getEnvValue("MT_DB_PASSWORD", cfg.DBPassword)
	cfg.DBPort = getEnvValue("MT_DB_PORT", cfg.DBPort)

	return cfg
}

func getEnvValue(envName string, def string) string {
	if val, exists := os.LookupEnv(envName); exists {
		return val
	}
	return def
}
