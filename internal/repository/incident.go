package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"log"
	"time"
)

type IncidentRepo struct {
	db *sqlx.DB
}

func NewIncidentPostgres(db *sqlx.DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (r *IncidentRepo) CreateIncident(req models.IncidentRequest) error {
	query := `
		INSERT INTO incidents 
		(type, description, location, radius_meters, is_active, created_at, updated_at)
		VALUES (
		$1, 
		$2, 
		ST_MakePoint($3, $4)::geography, 
		$5,
		$6,
		$7,
		$8)
		RETURNING id
	`

	now := time.Now()

	var id int64
	err := r.db.QueryRowx(
		query,
		req.Type,
		req.Description,
		req.Longitude,
		req.Latitude,
		req.RadiusMeters,
		req.Active,
		now,
		now,
	).Scan(&id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *IncidentRepo) GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error) {
	query := `
		SELECT
			type,
			description,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			radius_meters,
			is_active,
			created_at,
			updated_at
		FROM incidents
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []models.IncidentResponse

	for rows.Next() {
		var inc models.IncidentResponse

		if err := rows.Scan(
			&inc.Type,
			&inc.Description,
			&inc.Latitude,
			&inc.Longitude,
			&inc.RadiusMeters,
			&inc.Active,
			&inc.CreatedAt,
			&inc.UpdatedAt,
		); err != nil {
			return nil, err
		}

		resp = append(resp, inc)
	}

	return resp, rows.Err()
}

func (r *IncidentRepo) GetIncidentById(id int) (models.IncidentResponse, error) {
	query := `
		SELECT
			type,
			description,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			radius_meters,
			is_active,
			created_at,
			updated_at
		FROM incidents
		WHERE id = $1
	`

	var inc models.IncidentResponse
	err := r.db.Get(&inc, query, id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.IncidentResponse{}, sql.ErrNoRows
		}
		return models.IncidentResponse{}, err
	}

	return inc, nil
}

func (r *IncidentRepo) UpdateIncident(
	id int,
	req models.IncidentRequest,
) (models.IncidentResponse, error) {

	query := `
		UPDATE incidents
		SET
			type = $1,
			description = $2,
			location = ST_MakePoint($3, $4)::geography,
			radius_meters = $5,
			is_active = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING
			type,
			description,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			radius_meters,
			is_active,
			created_at,
			updated_at
	`

	var incident models.IncidentResponse

	err := r.db.Get(
		&incident,
		query,
		req.Type,
		req.Description,
		req.Latitude,
		req.Longitude,
		req.RadiusMeters,
		req.Active,
		id,
	)

	if err != nil {
		return models.IncidentResponse{}, err
	}

	return incident, nil
}

func (r *IncidentRepo) DeleteIncident(id int) error {
	query := `
		UPDATE incidents
		SET 
			is_active = false,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *IncidentRepo) GetDangerStats(
	ctx context.Context,
	window time.Duration,
) (int64, error) {

	const query = `
        SELECT COUNT(DISTINCT user_id)
        FROM location_checks
        WHERE
            has_danger = TRUE
            AND created_at >= now() - $1::interval
    `

	interval := fmt.Sprintf("%d seconds", int(window.Seconds()))

	var count int64
	err := r.db.QueryRowContext(ctx, query, interval).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *IncidentRepo) GetActiveIncidents(ctx context.Context) ([]models.IncidentResponse, error) {
	const query = `
        SELECT id, title, status, created_at
        FROM incidents
        WHERE status = 'active'
        ORDER BY created_at DESC
    `
	var incidents []models.IncidentResponse
	err := r.db.SelectContext(ctx, &incidents, query)
	return incidents, err
}
