package service

import (
	"context"
	"fmt"

	"github.com/tyler-garrett/blind-service/internal/domain"
	"github.com/tyler-garrett/blind-service/internal/repository"
)

type BlindService struct {
	repo *repository.BlindRepository
}

func NewBlindService(repo *repository.BlindRepository) *BlindService {
	return &BlindService{
		repo: repo,
	}
}

func (s *BlindService) GetAll(ctx context.Context) ([]domain.Blind, error) {
	blinds, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all blinds: %w", err)
	}
	return blinds, nil
}

func (s *BlindService) GetById(ctx context.Context, id string) (*domain.Blind, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	blind, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blind: %w", err)
	}
	return blind, nil
}

func (s *BlindService) GetInRadius(ctx context.Context, req domain.BlindsInRadiusRequest) ([]domain.Blind, error) {
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}
	if req.RadiusMeters <= 0 {
		return nil, fmt.Errorf("radius must be greater than 0")
	}

	blinds, err := s.repo.GetInRadius(ctx, req.Latitude, req.Longitude, req.RadiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to get blinds in radius: %w", err)
	}
	return blinds, nil
}

func (s *BlindService) Create(ctx context.Context, req domain.CreateBlindRequest) (*domain.Blind, error) {
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}
	if req.Id == "" {
		return nil, fmt.Errorf("blind Id is required")
	}

	blind := &domain.Blind{
		Id:        req.Id,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := s.repo.Create(ctx, blind); err != nil {
		return nil, fmt.Errorf("failed to create blind: %w", err)
	}

	return blind, nil
}

func (s *BlindService) Update(ctx context.Context, id string, req domain.UpdateBlindRequest) (*domain.Blind, error) {
	if id == "" {
		return nil, fmt.Errorf("blind Id is required")
	}

	// Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	blind := &domain.Blind{
		Id:        id,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := s.repo.Update(ctx, id, blind); err != nil {
		return nil, fmt.Errorf("failed to update blind: %w", err)
	}

	// Fetch updated blind
	updated, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated blind: %w", err)
	}

	return updated, nil
}

func (s *BlindService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("blind Id is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete blind: %w", err)
	}
	return nil
}

func (s *BlindService) GetStats(ctx context.Context) (*domain.BlindStats, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}
