package model

import "time"

type Sensor struct {
	ID          ID
	Name        SensorName
	Type        SensorType
	Location    Location
	Unit        TemperatureUnit
	Status      SensorStatus
	LastUpdated time.Time
	CreatedAt   time.Time
}
