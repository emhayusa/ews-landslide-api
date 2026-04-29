package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	RedisHost    string
	RedisPort    string
	KeycloakURL  string
	Realm        string
	AuthMethod   string // "keycloak" or "local"
	JWTSecret    string
	MQTTHost     string
	MQTTPort     string
	MQTTUser     string
	MQTTPass     string
	MQTTClientID string
	MQTTTopic    string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPass:       getEnv("DB_PASS", "password"),
		DBName:       getEnv("DB_NAME", "big_devops_db"),
		RedisHost:    getEnv("REDIS_HOST", "localhost"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		KeycloakURL:  getEnv("KEYCLOAK_URL", "http://keycloak:8080"),
		Realm:        getEnv("KEYCLOAK_REALM", "big-devops"),
		AuthMethod:   getEnv("AUTH_METHOD", "local"),
		JWTSecret:    getEnv("JWT_SECRET", "very-secret-key"),
		MQTTHost:     getEnv("MQTT_HOST", "localhost"),
		MQTTPort:     getEnv("MQTT_PORT", "1883"),
		MQTTUser:     getEnv("MQTT_USERNAME", "hi-target"),
		MQTTPass:     getEnv("MQTT_PASSWORD", "h!-t4rg3t"),
		MQTTClientID: getEnv("MQTT_CLIENT_ID", "hi-target-api-stream"),
		MQTTTopic:    getEnv("MQTT_TOPIC", "gnss/hitarget"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
