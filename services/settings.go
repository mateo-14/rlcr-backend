package services

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Nico-14/rlcr-backend/db"
	"github.com/Nico-14/rlcr-backend/models"
)

type ISettingsService interface {
	Get(ctx context.Context) (*models.Settings, error)
	Update(ctx context.Context, settings *models.Settings) error
}

type SettingsService struct {
	db *db.Client
}

func NewSettingsService(db *db.Client) *SettingsService {
	return &SettingsService{db}
}

func (s *SettingsService) Get(ctx context.Context) (*models.Settings, error) {
	var settings models.Settings

	dsnap, err := s.db.Collection("settings").Doc("default").Get(ctx)
	if err != nil {
		return nil, err
	}

	err = dsnap.DataTo(&settings)
	return &settings, err
}

func (s *SettingsService) Update(ctx context.Context, settings *models.Settings) error {
	settm, _ := settings.ToMap()
	_, err := s.db.Collection("settings").Doc("default").Set(ctx, settm, firestore.MergeAll)
	return err
}
