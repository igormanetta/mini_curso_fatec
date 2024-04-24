package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// AppConfig contém as configurações da aplicação.
type AppConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	ImagesFolderPath string
}

// LoadConfig carrega as configurações da aplicação a partir de variáveis de ambiente.
func LoadConfig() AppConfig {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := parseIntEnvOrDefault("DB_PORT", 5432)
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASS", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "postgres")
	dbPath := getEnvOrDefault("FOLDER_PATH", `.\basedir\`)

	config := AppConfig{
		Host:             dbHost,
		Port:             dbPort,
		User:             dbUser,
		Password:         dbPass,
		Database:         dbName,
		ImagesFolderPath: dbPath,
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Variável de ambiente %s não encontrada. Usando valor padrão: %s\n", key, defaultValue)
		return defaultValue
	}
	return value
}

func parseIntEnvOrDefault(key string, defaultValue int) int {
	valueStr := getEnvOrDefault(key, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Erro ao converter variável de ambiente %s para inteiro. Usando valor padrão: %d\n", key, defaultValue)
		return defaultValue
	}
	return value
}
