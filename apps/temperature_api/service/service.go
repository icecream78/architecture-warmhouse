package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
)

type sensorRepository interface {
	GetSensorDetailsByLocation(ctx context.Context, location model.Location) (model.Sensor, error)
	GetSensorDetailsByID(ctx context.Context, id model.ID) (model.Sensor, error)
}

type sensorsClient interface {
	GetTemperatureFromLocation(ctx context.Context, location model.Location) (RealtimeDataForTemperatureByLocation, error)
}

type Service struct {
	repository    sensorRepository
	sensorsClient sensorsClient
}

func New(
	repository sensorRepository,
	sensorsClient sensorsClient,
) (*Service, error) {
	return &Service{
		repository:    repository,
		sensorsClient: sensorsClient,
	}, nil
}

func (s *Service) GetTemperatureByLocation(ctx context.Context, location model.Location) (TemperatureDataBySensor, error) {
	sensorDetails, err := s.repository.GetSensorDetailsByLocation(ctx, location)
	if err != nil {
		return TemperatureDataBySensor{}, errors.Wrap(err, "repo get sensor details")
	}

	if sensorDetails.Type != model.TemperatureSensorType {
		return TemperatureDataBySensor{}, errors.Wrap(err, "request temperature not from needed sensor")
	}

	realtimeData, err := s.sensorsClient.GetTemperatureFromLocation(ctx, location)
	if err != nil {
		return TemperatureDataBySensor{}, errors.Wrap(err, "request sensor temperature")
	}

	return TemperatureDataBySensor{
		SensorDetails: SensorDetailsForTemperatureByLocation{
			ID:          sensorDetails.ID,
			Type:        sensorDetails.Type,
			Unit:        sensorDetails.Unit,
			Status:      sensorDetails.Status,
			Location:    sensorDetails.Location,
			Description: "wtf to insert here???",
		},
		RealtimeData: RealtimeDataForTemperatureByLocation{
			Value:     realtimeData.Value,
			Timestamp: realtimeData.Timestamp,
		},
	}, nil
}

func (s *Service) GetTemperatureBySensorID(ctx context.Context, sensorID model.ID) (TemperatureDataBySensor, error) {
	sensorDetails, err := s.repository.GetSensorDetailsByID(ctx, sensorID)
	if err != nil {
		return TemperatureDataBySensor{}, errors.Wrap(err, "repo get sensor details")
	}

	if sensorDetails.Type != model.TemperatureSensorType {
		return TemperatureDataBySensor{}, errors.Wrap(err, "request temperature not from needed sensor")
	}

	realtimeData, err := s.sensorsClient.GetTemperatureFromLocation(ctx, sensorDetails.Location)
	if err != nil {
		return TemperatureDataBySensor{}, errors.Wrap(err, "request sensor temperature")
	}

	return TemperatureDataBySensor{
		SensorDetails: SensorDetailsForTemperatureByLocation{
			ID:          sensorDetails.ID,
			Type:        sensorDetails.Type,
			Unit:        sensorDetails.Unit,
			Status:      sensorDetails.Status,
			Location:    sensorDetails.Location,
			Description: "wtf to insert here???",
		},
		RealtimeData: RealtimeDataForTemperatureByLocation{
			Value:     realtimeData.Value,
			Timestamp: realtimeData.Timestamp,
		},
	}, nil
}
