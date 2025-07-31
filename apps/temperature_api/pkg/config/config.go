package config

import "os"

type Config struct {
	Environment string
	HttpServer  HttpServer
	Database    Database
}

type HttpServer struct {
	URL string
}

type Database struct {
	URL string
}

func Get() (*Config, error) {
	environment := "production"
	if url := os.Getenv("ENV"); url != "" {
		environment = url
	}

	httpListenURL := "0.0.0.0:8081"
	if url := os.Getenv("LISTEN_ADDRESS"); url != "" {
		httpListenURL = url
	}

	databaseURL := "postgres://postgres:postgres@localhost:5432/smarthome"
	if url := os.Getenv("DATABASE_URL"); url != "" {
		databaseURL = url
	}

	return &Config{
		Environment: environment,
		HttpServer: HttpServer{
			URL: httpListenURL,
		},
		Database: Database{
			URL: databaseURL,
		},
	}, nil
}
