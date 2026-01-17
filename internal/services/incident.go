package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"log"
	"time"
)

type IncidentService struct {
	repo         repository.Incident
	cache        repository.IncidentCache
	windowMin    int
	webhookQueue WebhookQueue
}

func NewIncidentService(repo repository.Incident, windowMin int, webhookQueue WebhookQueue) *IncidentService {
	return &IncidentService{repo: repo, windowMin: windowMin, webhookQueue: webhookQueue}
}

var ErrIncidentAlreadyExists = errors.New("incident already exists")

func (s *IncidentService) CheckLocation(ctx context.Context, checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {

	nearbyResp, err := s.repo.CheckLocation(checkReq)

	if err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
	}

	err = s.repo.SaveCheck(checkReq.UserID, checkReq.Lat, checkReq.Lon, true)
	if err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
	}

	if nearbyResp.Danger {
		job := models.WebhookPayload{
			Event:     "danger_detected",
			UserID:    checkReq.UserID,
			Lat:       checkReq.Lat,
			Lon:       checkReq.Lon,
			Incidents: nearbyResp.Incidents,
			CheckedAt: time.Now().UTC(),
		}

		_ = s.webhookQueue.Enqueue(ctx, job)
	}

	return nearbyResp, nil
}

func (s *IncidentService) CreateIncident(req models.IncidentRequest) error {
	err := s.repo.CreateIncident(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentService) GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error) {
	incidents, err := s.repo.GetAllIncidents(limit, offset)
	if err != nil {
		return []models.IncidentResponse{}, err
	}

	return incidents, nil
}

func (s *IncidentService) GetIncidentById(id int) (models.IncidentResponse, error) {
	incident, err := s.repo.GetIncidentById(id)
	if err != nil {
		return models.IncidentResponse{}, err
	}

	return incident, nil
}

func (s *IncidentService) UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error) {
	_, err := s.repo.GetIncidentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return models.IncidentResponse{}, err
		}
		log.Println(err)
		return models.IncidentResponse{}, err
	}

	updateIncident, err := s.repo.UpdateIncident(id, req)
	if err != nil {
		log.Println(err)
		return models.IncidentResponse{}, err
	}

	return updateIncident, nil
}

func (s *IncidentService) DeleteIncident(id int) error {
	err := s.repo.DeleteIncident(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentService) GetIncidentStats(ctx context.Context) (models.IncidentStatsResponse, error) {

	window := time.Duration(s.windowMin) * time.Minute

	count, err := s.repo.GetDangerStats(ctx, window)
	if err != nil {
		return models.IncidentStatsResponse{}, err
	}

	return models.IncidentStatsResponse{
		UserCount:    count,
		WindowMinute: s.windowMin,
	}, nil
}

func (s *IncidentService) GetActiveIncidents(ctx context.Context) ([]models.IncidentResponse, error) {
	if cached, err := s.cache.GetActive(ctx); err == nil && cached != nil {
		return cached, nil
	}

	incidents, err := s.repo.GetActiveIncidents(ctx)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetActive(ctx, incidents, 30*time.Second)

	return incidents, nil
}
