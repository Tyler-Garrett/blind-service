package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/tyler-garrett/blind-service/internal/domain"
	"github.com/tyler-garrett/blind-service/pkg/distance"
)

type BlindRepository struct {
	client *aztables.Client
}

type BlindEntity struct {
	aztables.Entity
	Id        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func NewBlindRepository(client *aztables.Client) *BlindRepository {
	return &BlindRepository{
		client: client,
	}
}

func (r *BlindRepository) GetAll(ctx context.Context) ([]domain.Blind, error) {
	var blinds []domain.Blind

	filter := "PartitionKey eq 'VA'"
	pager := r.client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: &filter,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %w", err)
		}

		for _, entity := range page.Entities {
			var blindEntity BlindEntity
			if err := json.Unmarshal(entity, &blindEntity); err != nil {
				// Log error but continue processing other entities
				continue
			}

			blind := r.entityToDomain(&blindEntity)
			blinds = append(blinds, blind)
		}
	}

	return blinds, nil
}

func (r *BlindRepository) GetById(ctx context.Context, id string) (*domain.Blind, error) {
	response, err := r.client.GetEntity(ctx, "VA", id, nil)
	if err != nil {
		return nil, domain.ErrBlindNotFound
	}

	var blindEntity BlindEntity
	if err := json.Unmarshal(response.Value, &blindEntity); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity: %w", err)
	}

	blind := r.entityToDomain(&blindEntity)
	return &blind, nil
}

func (r *BlindRepository) GetInRadius(ctx context.Context, lat, lon, radiusMeters float64) ([]domain.Blind, error) {
	allBlinds, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.Blind
	for _, blind := range allBlinds {
		dist := distance.Calculate(lat, lon, blind.Latitude, blind.Longitude)
		if dist <= radiusMeters {
			result = append(result, blind)
		}
	}

	return result, nil
}

func (r *BlindRepository) Create(ctx context.Context, blind *domain.Blind) error {
	existing, _ := r.GetById(ctx, blind.Id)
	if existing != nil {
		return domain.ErrBlindAlreadyExists
	}

	entity := r.domainToEntity(blind)

	marshaled, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}

	_, err = r.client.AddEntity(ctx, marshaled, nil)
	if err != nil {
		return fmt.Errorf("failed to create blind: %w", err)
	}

	return nil
}

func (r *BlindRepository) Update(ctx context.Context, id string, blind *domain.Blind) error {
	existing, err := r.GetById(ctx, id)
	if err != nil {
		return err
	}

	existing.Latitude = blind.Latitude
	existing.Longitude = blind.Longitude

	updatedAt := time.Now()
	existing.UpdatedAt = &updatedAt

	entity := r.domainToEntity(existing)

	marshaled, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}

	_, err = r.client.UpdateEntity(ctx, marshaled, &aztables.UpdateEntityOptions{
		UpdateMode: aztables.UpdateModeReplace,
	})
	if err != nil {
		return fmt.Errorf("failed to update blind: %w", err)
	}

	return nil
}

func (r *BlindRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteEntity(ctx, "VA", id, nil)
	if err != nil {
		return domain.ErrBlindNotFound
	}
	return nil
}

// GetStats returns statistics about all blinds
func (r *BlindRepository) GetStats(ctx context.Context) (*domain.BlindStats, error) {
	blinds, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(blinds) == 0 {
		return &domain.BlindStats{}, nil
	}

	stats := &domain.BlindStats{
		TotalBlinds:  len(blinds),
		MinLatitude:  blinds[0].Latitude,
		MaxLatitude:  blinds[0].Latitude,
		MinLongitude: blinds[0].Longitude,
		MaxLongitude: blinds[0].Longitude,
	}

	for _, blind := range blinds {
		if blind.Latitude < stats.MinLatitude {
			stats.MinLatitude = blind.Latitude
		}
		if blind.Latitude > stats.MaxLatitude {
			stats.MaxLatitude = blind.Latitude
		}
		if blind.Longitude < stats.MinLongitude {
			stats.MinLongitude = blind.Longitude
		}
		if blind.Longitude > stats.MaxLongitude {
			stats.MaxLongitude = blind.Longitude
		}
	}

	return stats, nil
}

func (r *BlindRepository) entityToDomain(entity *BlindEntity) domain.Blind {
	blind := domain.Blind{
		Id:        entity.Id,
		Latitude:  entity.Latitude,
		Longitude: entity.Longitude,
	}

	if createdAt, err := time.Parse(time.RFC3339, entity.CreatedAt); err == nil {
		blind.CreatedAt = createdAt
	}

	if entity.UpdatedAt != "" {
		if updatedAt, err := time.Parse(time.RFC3339, entity.UpdatedAt); err == nil {
			blind.UpdatedAt = &updatedAt
		}
	}

	return blind
}

func (r *BlindRepository) domainToEntity(blind *domain.Blind) *BlindEntity {
	entity := &BlindEntity{
		Entity: aztables.Entity{
			PartitionKey: "VA",
			RowKey:       blind.Id,
		},
		Id:        blind.Id,
		Latitude:  blind.Latitude,
		Longitude: blind.Longitude,
		CreatedAt: blind.CreatedAt.Format(time.RFC3339),
	}

	if blind.UpdatedAt != nil {
		entity.UpdatedAt = blind.UpdatedAt.Format(time.RFC3339)
	}

	return entity
}
