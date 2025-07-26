package service

import (
	"time"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
)

type TemperatureByLocation struct {
	RealtimeData  RealtimeDataForTemperatureByLocation
	SensorDetails SensorDetailsForTemperatureByLocation
}

type RealtimeDataForTemperatureByLocation struct {
	Value     model.SensorValue
	Timestamp time.Time
}

type SensorDetailsForTemperatureByLocation struct {
	ID          model.ID
	Type        model.SensorType
	Unit        model.TemperatureUnit
	Status      model.SensorStatus
	Location    model.Location
	Description string
}
