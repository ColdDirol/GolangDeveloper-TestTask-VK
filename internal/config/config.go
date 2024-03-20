package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ENV             string `json:"env"`
	IntervalSeconds int    `json:"interval_seconds"`
	MaxRequests     int    `json:"max_requests"`
	Redis           Redis
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func InitConfig(configPath string) *Config {
	// Читаем файл
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("ошибка чтения файла конфигурации: %v", err)
	}

	// Парсим JSON
	var config Config
	if err := json.Unmarshal(fileContent, &config); err != nil {
		log.Fatalf("ошибка парсинга JSON: %v", err)
	}
	if os.Getenv("REDIS_HOST") != "" {
		config.Redis.Host = os.Getenv("REDIS_HOST")
	}

	return &config
}
