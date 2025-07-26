package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/client/sensor"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/handler"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/pkg/config"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/repository"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/service"
)

type exitCode int

const (
	succesExitCode exitCode = 0
	errorExitCode  exitCode = 1
)

func main() {
	os.Exit(int(run()))
}

type sensorService interface {
	GetTemperatureByLocation(ctx context.Context, location model.Location) (service.TemperatureByLocation, error)
}

func run() exitCode {
	appConfig, err := config.Get()
	if err != nil {
		slog.Error("during get config", slog.String("error", err.Error()))
		return errorExitCode
	}

	sensorsClient, err := sensor.New()
	if err != nil {
		slog.Error("during init sensor client", slog.String("error", err.Error()))
		return errorExitCode
	}

	sensorRepository, err := repository.New(appConfig.Database.URL)
	if err != nil {
		slog.Error("during init repository", slog.String("error", err.Error()))
		return errorExitCode
	}

	sensorService, err := service.New(sensorRepository, sensorsClient)
	if err != nil {
		slog.Error("during init sensor service", slog.String("error", err.Error()))
		return errorExitCode
	}

	e := echo.New()

	if err := handler.RegisterRouteHandlers(sensorService, e); err != nil {
		slog.Error("during register http handlers", slog.String("error", err.Error()))
		return errorExitCode
	}

	e.Logger.Fatal(e.Start(appConfig.HttpServer.URL))

	return succesExitCode
}
