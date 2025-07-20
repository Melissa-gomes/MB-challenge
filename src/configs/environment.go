package configs

import (
	"os"
	"strconv"
)

type Env struct {
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     int
	SSLMODE           string
}

func LoadEnv() Env {
	envPort := os.Getenv("POSTGRES_PORT")
	port, _ := strconv.Atoi(envPort)
	return Env{
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     port,
		SSLMODE:           os.Getenv("SSLMODE"),
	}
}
