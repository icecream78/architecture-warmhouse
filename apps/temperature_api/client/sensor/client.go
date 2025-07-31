package sensor

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/service"
)

const (
	temperatureMinValue float64 = 10
	temperatureMaxValue float64 = 30
)

type Client struct{}

func New() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) GetTemperatureFromLocation(ctx context.Context, location model.Location) (service.RealtimeDataForTemperatureByLocation, error) {
	return service.RealtimeDataForTemperatureByLocation{
		Timestamp: time.Now().UTC(),
		Value:     model.SensorValue(generateRandomValue(temperatureMinValue, temperatureMaxValue)),
	}, nil
}

func generateRandomValue(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
