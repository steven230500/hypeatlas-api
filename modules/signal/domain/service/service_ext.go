package service

import (
	"context"
	"errors"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

func (s *svc) ListChampions(ctx context.Context, version string) ([]entities.Champion, error) {
	if s.riotSvc == nil {
		return nil, errors.New("riot service not available")
	}

	// Si no se especifica versión, obtener la última
	if version == "" {
		return nil, errors.New("version required")
	}

	riotChamps, err := s.riotSvc.GetChampions(ctx, version)
	if err != nil {
		return nil, err
	}

	var champions []entities.Champion
	for _, c := range riotChamps.Data {
		champions = append(champions, entities.Champion{
			ID:    c.ID,
			Key:   c.Key,
			Name:  c.Name,
			Title: c.Title,
		})
	}

	return champions, nil
}

func (s *svc) ListRegions(ctx context.Context) ([]string, error) {
	if s.riotSvc == nil {
		return nil, errors.New("riot service not available")
	}
	return s.riotSvc.GetRegions(ctx)
}
