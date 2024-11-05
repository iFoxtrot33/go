package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type RecoveryConfig struct {
	SMTP     SMTPConfig `json:"smtp"`
	DataFile string     `json:"data_file"`
}

type Config struct {
	Recovery RecoveryConfig `json:"recovery"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Recovery: RecoveryConfig{
			SMTP: SMTPConfig{
				Host:     "live.smtp.mailtrap.io",
				Port:     587,
				Username: "api",
				Password: "dd148f4f71dcc620d648a25ec81e241b",
				From:     "noreply@demomailtrap.com",
			},
			DataFile: "data.json",
		},
	}
}

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Load() *Config {
	return &Config{
		Recovery: RecoveryConfig{
			SMTP: SMTPConfig{
				Host:     getEnv("SMTP_HOST", "live.smtp.mailtrap.io"),
				Port:     getEnvInt("SMTP_PORT", 587),
				Username: getEnv("SMTP_USERNAME", "api"),
				Password: getEnv("SMTP_PASSWORD", "dd148f4f71dcc620d648a25ec81e241b"),
				From:     getEnv("SMTP_FROM", "noreply@demomailtrap.com"),
			},
			DataFile: getEnv("DATA_FILE", "data.json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
