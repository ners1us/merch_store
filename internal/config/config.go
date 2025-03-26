package config

import "os"

type Config struct {
	DbUrl     string
	JWTSecret string
	Port      string
}

func InitConfig() *Config {
	return &Config{
		DbUrl:     getEnv("DB_URL"),
		JWTSecret: getEnv("JWT_SECRET"),
		Port:      getEnv("PORT"),
	}
}

func getEnv(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}
