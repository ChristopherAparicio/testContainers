package config

import (
	"log"

	env "github.com/Netflix/go-env"
)

type Config struct {
	Server Server

	// Only Postgresql for now
	Database PostgresqlDatabase `json:"database"`
}

type Server struct {
	Port int `json:"port" env:"SERVER_PORT,default=8080"`
}

// Use Netflix go env
type PostgresqlDatabase struct {
	Host     string `json:"host" env:"DATABASE_HOST,default=localhost"`
	Port     int    `json:"port" env:"DATABASE_PORT,default=5432"`
	User     string `json:"user" env:"DATABASE_USER,default=postgres"`
	Password string `json:"password" env:"DATABASE_PASSWORD,default=postgres"`
	Database string `json:"database" env:"DATABASE_NAME,default=tinyurl"`
}

func ParseConfig() Config {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
