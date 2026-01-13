package repository

import (
	"database/sql"
	"errors"
	"time"
	"github.com/jmoiron/sqlx"
	"context"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

type IncidentRepo struct {
	// db *sql.DB
	db *sqlx.DB
}

func NewIncidentPostgres(db *sqlx.DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

// func (r *IncidentRepo) GetIncident(req models.IncidentRequest) (int, error) {
// 	query := `
// 		SELECT id, title, description, status, active, created_by, updated_by, created_at, updated_at, closed_at
// 		FROM incidents
// 		WHERE title = , description, status, active, created_by, updated_by, created_at, updated_at, closed_at
// 	`

// 	var inc models.Incident
// 	err := r.db.Get(&inc, query, id)
// 	if err != nil {
// 		if errors.Is(err, sqlx.ErrNotFound) {
// 			return 0, sqlx.ErrNotFound
// 		}
// 		return 0, err
// 	}

// 	return &inc, nil
// }



// func (r *IncidentRepo) FindNearby(lat, lon) ([]domain.NearbyIncident, error) {
func (r *IncidentRepo) CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {
	const query = `
        WITH user_point AS (
            SELECT ST_MakePoint($1, $2)::geography AS geom
        )
        SELECT
            i.id,
            i.type,
            ST_Distance(i.location, up.geom) AS distance_meters
        FROM incidents i, user_point up
        WHERE
            i.is_active = TRUE
            AND (i.starts_at IS NULL OR i.starts_at <= now())
            AND (i.ends_at IS NULL OR i.ends_at >= now())
            AND ST_DWithin(
                i.location,
                up.geom,
                i.radius_meters
            )
        ORDER BY distance_meters;
    `

	rows, err := r.db.Query(
        query,
        checkReq.Lon, // ⚠️ сначала lon
        checkReq.Lat, // потом lat
    )

	if err != nil {
        return models.LocationCheckResponse{}, err
    }
    defer rows.Close()

    resp := models.LocationCheckResponse{}

    for rows.Next() {
        var inc models.NearbyIncidentResponse

        if err := rows.Scan(
            &inc.ID,
            &inc.Type,
            &inc.DistanceMeters,
        ); err != nil {
            return models.LocationCheckResponse{}, err
        }

        resp.Incidents = append(resp.Incidents, inc)
    }

    if err := rows.Err(); err != nil {
        return models.LocationCheckResponse{}, err
    }

    resp.Danger = len(resp.Incidents) > 0

    return resp, nil
}

func (r *IncidentRepo) SaveCheck(userID string, lat, lon float64, hasDanger bool,) error {
	const query = `
        INSERT INTO location_checks (user_id, location, has_danger)
        VALUES ($1, ST_MakePoint($2, $3)::geography, $4)
    `

    _, err := r.db.ExecContext(
		context.Background(),
        query,
        userID,
        lon, // порядок!
        lat,
        hasDanger,
    )

    return err
}

func (r *IncidentRepo) CreateIncident(req models.IncidentRequest) error {
	query := `
		INSERT INTO incidents 
		(title, description, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()

	var id int64
	err := r.db.QueryRowx(
		query,
		req.Title,
		req.Description,
		req.Active,
		now,
		now,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (r *IncidentRepo) GetIncidentById(id int) (models.Incident, error) {
	query := `
		SELECT title, description, active, created_at, updated_at
		FROM incidents
		WHERE id = $1
	`
	
	var inc models.Incident
	err := r.db.Get(&inc, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Incident{}, sql.ErrNoRows
		}
		return models.Incident{}, err
	}

	return inc, nil
}

// func (r *IncidentRepo) UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error) {
// 	query := `
// 		UPDATE incidents
// 		SET
// 			title = $1,
// 			description = $2,
// 			active = $4,
// 			updated_at = NOW()
// 		WHERE id = $5
// 		RETURNING
// 			id,
// 			title,
// 			description,
// 			active,
// 			created_at,
// 			updated_at,
// 	`
// 	var title string
// 	if req.Title != nil {
// 		title = *req.Title
// 	}

// 	var description string
// 	if req.Description != nil {
// 		description = *req.Description
// 	}

// 	var active bool
// 	if req.Active != nil {
// 		active = *req.Active
// 	}

// 	var incident models.IncidentResponse

// 	err := r.db.Get(
// 		&incident,
// 		query,
// 		title,
// 		description,
// 		active,
// 		id,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return models.IncidentResponse{}, sql.ErrNoRows
// 		}
// 		return models.IncidentResponse{}, err
// 	}

// 	return incident, nil

// }

func (r *IncidentRepo) UpdateIncident(
	id int,
	req models.IncidentRequest,
) (models.Incident, error) {

	query := `
		UPDATE incidents
		SET
			title = $1,
			description = $2,
			active = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING
			id,
			title,
			description,
			active,
			created_at,
			updated_at
	`

	var incident models.Incident

	err := r.db.Get(
		&incident,
		query,
		req.Title,
		req.Description,
		req.Active,
		id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Incident{}, sql.ErrNoRows
		}
		return models.Incident{}, err
	}

	return incident, nil
}

func (r *IncidentRepo) DeleteIncident(id int) error {
	query := `
		UPDATE incidents
		SET 
			active = false,
			updated_at = NOW()
		WHERE id = $1 AND active = true
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}