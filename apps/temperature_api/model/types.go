package model

import (
	"errors"
	"strings"
)

type SensorValue float64

type ID string

func IDFromRawString(raw string) (ID, error) {
	return ID(raw), nil
}

func (id ID) ToString() string {
	return string(id)
}

type SensorName string

func SensorNameFromRawString(raw string) (SensorName, error) {
	if strings.Contains(raw, ".") {
		return "", errors.New("sensor name can not include dots")
	}

	// подчистим от лишних пробелов
	raw = strings.TrimSpace(raw)

	return SensorName(raw), nil
}

func (sn SensorName) ToString() string {
	return string(sn)
}

type Location string

func LocationFromRawString(raw string) (Location, error) {
	if strings.Contains(raw, ".") {
		return "", errors.New("location can not include dots")
	}

	// подчистим от лишних пробелов
	raw = strings.TrimSpace(raw)

	return Location(raw), nil
}

func (ln Location) ToString() string {
	return string(ln)
}

type TemperatureUnit string

const (
	CelsiusTemperatureUnit    TemperatureUnit = "celsius"
	FahrenheitTemperatureUnit TemperatureUnit = "аahrenheit"
)

func TemperatureUnitFromRawString(raw string) (TemperatureUnit, error) {
	switch TemperatureUnit(raw) {
	case CelsiusTemperatureUnit, FahrenheitTemperatureUnit:
		return TemperatureUnit(raw), nil
	default:
		return "", errors.New("unknown temperature unit")
	}
}

func (tu TemperatureUnit) ToString() string {
	return string(tu)
}

func (tu TemperatureUnit) ToShortString() string {
	switch tu {
	case CelsiusTemperatureUnit:
		return "°C"
	case FahrenheitTemperatureUnit:
		return "°F"
	}

	return "unknown"
}

type SensorType string

func (st SensorType) ToString() string {
	return string(st)
}

const (
	TemperatureSensorType SensorType = "temperature"
	// все оставшиеся типы сенсоров по-умолчанию уносятся в категорию сервисных
	ServiceSensorType SensorType = "service"
)

type SensorStatus string

func (ss SensorStatus) ToString() string {
	return string(ss)
}
