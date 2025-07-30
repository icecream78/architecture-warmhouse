package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/service"
)

type sensorService interface {
	GetTemperatureByLocation(ctx context.Context, location model.Location) (service.TemperatureDataBySensor, error)
	GetTemperatureBySensorID(ctx context.Context, id model.ID) (service.TemperatureDataBySensor, error)
}

type sensorRepository interface {
	Ping(ctx context.Context) error
}

type handler struct {
	sensorService    sensorService
	sensorRepository sensorRepository
}

func RegisterRouteHandlers(
	sensorService sensorService,
	sensorRepository sensorRepository,
	e *echo.Echo,
) error {
	h := handler{
		sensorService:    sensorService,
		sensorRepository: sensorRepository,
	}

	e.GET("/temperature", h.GetTemperatureByLocation)
	e.GET("/temperature/:id", h.GetTemperatureBySensorID)

	e.GET("/health", h.HealthCheck)

	return nil
}

func (h *handler) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()

	status := "ok"
	isPgAvailable := h.sensorRepository.Ping(ctx) == nil
	if !isPgAvailable {
		status = "failure"
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": status,
		"resources": map[string]any{
			"postgres": isPgAvailable,
		},
	})

	return nil
}

func (h *handler) GetTemperatureBySensorID(c echo.Context) error {
	ctx := c.Request().Context()

	sensorID, err := model.IDFromRawString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Not value sensorID",
		})

		return nil
	}

	temperature, err := h.sensorService.GetTemperatureBySensorID(ctx, sensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})

		return nil
	}

	c.JSON(http.StatusOK, getTemperatureByLocationDTO{
		Value:       float64(temperature.RealtimeData.Value),
		Unit:        temperature.SensorDetails.Unit.ToShortString(),
		Timestamp:   temperature.RealtimeData.Timestamp,
		Location:    temperature.SensorDetails.Location.ToString(),
		Status:      temperature.SensorDetails.Status.ToString(),
		SensorID:    temperature.SensorDetails.ID.ToString(),
		SensorType:  temperature.SensorDetails.Type.ToString(),
		Description: temperature.SensorDetails.Description,
	})

	return nil
}

func (h *handler) GetTemperatureByLocation(c echo.Context) error {
	ctx := c.Request().Context()

	location, err := model.LocationFromRawString(c.QueryParam("location"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Not supported location name. Try to fix it and repeat",
		})

		return nil
	}

	temperature, err := h.sensorService.GetTemperatureByLocation(ctx, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})

		return nil
	}

	c.JSON(http.StatusOK, getTemperatureByLocationDTO{
		Value:       float64(temperature.RealtimeData.Value),
		Unit:        temperature.SensorDetails.Unit.ToShortString(),
		Timestamp:   temperature.RealtimeData.Timestamp,
		Location:    temperature.SensorDetails.Location.ToString(),
		Status:      temperature.SensorDetails.Status.ToString(),
		SensorID:    temperature.SensorDetails.ID.ToString(),
		SensorType:  temperature.SensorDetails.Type.ToString(),
		Description: temperature.SensorDetails.Description,
	})

	return nil
}
