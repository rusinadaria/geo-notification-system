package services

import (
	"context"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"log"
	"time"
)

type locationCheckService struct {
	repo         repository.LocationCheck
	webhookQueue WebhookQueue
}

func NewLocationCheckService(repo repository.LocationCheck, webhookQueue WebhookQueue) *locationCheckService {
	return &locationCheckService{repo: repo, webhookQueue: webhookQueue}
}

func (l *locationCheckService) CheckLocation(ctx context.Context, checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {

	nearbyResp, err := l.repo.CheckLocation(checkReq)

	if err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
	}

	err = l.repo.SaveCheck(checkReq.UserID, checkReq.Lat, checkReq.Lon, true)
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

		_ = l.webhookQueue.Enqueue(ctx, job)
	}

	return nearbyResp, nil
}
