package config

import "os"

type Config struct {
	HttpServer HttpServer
	Database   Database
}

type HttpServer struct {
	URL string
}

type Database struct {
	URL string
}

func Get() (*Config, error) {
	httpListenURL := "http://0.0.0.0:8080"
	if url := os.Getenv("LISTEN_ADDRESS"); url != "" {
		httpListenURL = url
	}

	databaseURL := "postgres://postgres:postgres@localhost:5432/smarthome"
	if url := os.Getenv("DATABASE_URL"); url != "" {
		databaseURL = url
	}

	return &Config{
		HttpServer: HttpServer{
			URL: httpListenURL,
		},
		Database: Database{
			URL: databaseURL,
		},
	}, nil
}
