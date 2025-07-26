package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/service"
)

type handler struct {
	sensorService *service.Service
}

func RegisterRouteHandlers(sensorService *service.Service, e *echo.Echo) error {
	h := handler{
		sensorService: sensorService,
	}

	e.GET("/temperature", h.GetTemperatureByLocation)

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
