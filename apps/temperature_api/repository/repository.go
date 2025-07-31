package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/icecream78/architecture-warmhouse/apps/temperature_api/model"
)

type Repository struct {
	Pool *pgxpool.Pool
}

// New creates a new DB instance
func New(connString string) (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Repository{Pool: pool}, nil
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.Pool.Ping(ctx)
}

// Close closes the database connection
func (r *Repository) Close() {
	if r.Pool != nil {
		r.Pool.Close()
	}
}

func (r *Repository) GetSensorDetailsByLocation(ctx context.Context, location model.Location) (model.Sensor, error) {
	query := `
		SELECT id, name, type, unit, status, last_updated, created_at
		FROM sensors
		WHERE location = $1
	`

	var s model.Sensor
	err := r.Pool.QueryRow(ctx, query, location.ToString()).Scan(
		&s.ID,
		&s.Name,
		&s.Type,
		&s.Unit,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
	)
	if err != nil {
		return model.Sensor{}, fmt.Errorf("error getting sensor by location: %w", err)
	}

	s.Location = location

	return s, nil
}

func (r *Repository) GetSensorDetailsByID(ctx context.Context, id model.ID) (model.Sensor, error) {
	query := `
		SELECT id, name, type, unit, status, last_updated, created_at, location
		FROM sensors
		WHERE id = $1
	`

	var s model.Sensor
	err := r.Pool.QueryRow(ctx, query, id.ToString()).Scan(
		&s.ID,
		&s.Name,
		&s.Type,
		&s.Unit,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
		&s.Location,
	)
	if err != nil {
		return model.Sensor{}, fmt.Errorf("error getting sensor by sensor id: %w", err)
	}

	return s, nil
}
