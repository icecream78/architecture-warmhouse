package sensor

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/service"
)

type Client struct{}

func New() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) GetTemperatureFromLocation(ctx context.Context, location model.Location) (service.RealtimeDataForTemperatureByLocation, error) {
	return service.RealtimeDataForTemperatureByLocation{
		Timestamp: time.Now().UTC(),
		Value:     model.SensorValue(rand.Float64()),
	}, nil
}
