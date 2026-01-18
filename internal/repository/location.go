package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"log"
)

type LocationCheckRepo struct {
	db *sqlx.DB
}

func NewLocationCheckPostgres(db *sqlx.DB) *LocationCheckRepo {
	return &LocationCheckRepo{db: db}
}

func (r *LocationCheckRepo) CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {
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
			AND ST_DWithin(
				i.location,
				up.geom,
				i.radius_meters
			)
		ORDER BY distance_meters;
    `

	rows, err := r.db.Query(
		query,
		checkReq.Lon,
		checkReq.Lat,
	)

	if err != nil {
		log.Println(err)
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
			log.Println(err)
			return models.LocationCheckResponse{}, err
		}

		resp.Incidents = append(resp.Incidents, inc)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
	}

	resp.Danger = len(resp.Incidents) > 0

	return resp, nil
}

func (r *LocationCheckRepo) SaveCheck(userID int, lat, lon float64, hasDanger bool) error {
	const query = `
        INSERT INTO location_checks (user_id, location, has_danger)
        VALUES ($1, ST_MakePoint($2, $3)::geography, $4)
    `

	_, err := r.db.ExecContext(
		context.Background(),
		query,
		userID,
		lon,
		lat,
		hasDanger,
	)

	return err
}
