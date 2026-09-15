package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/zone/dto"
	"github.com/epmp/backend/internal/modules/zone/entity"
	"github.com/epmp/backend/internal/modules/zone/repository"
)

// ZoneService implements the application layer for Zone.
type ZoneService struct {
	repo repository.ZoneRepository
}

// NewZoneService creates a new ZoneService.
func NewZoneService(repo repository.ZoneRepository) *ZoneService {
	return &ZoneService{repo: repo}
}

func (s *ZoneService) Create(ctx context.Context, orgID string, req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("zone service: create: organization ID required")
	}
	e := entity.NewZone()
	e.OrganizationId = orgID
	e.BuildingId = req.BuildingId
	e.Floor = req.Floor
	e.Name = req.Name

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("zone service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ZoneService) GetByID(ctx context.Context, id, orgID string) (*dto.ZoneResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("zone service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *ZoneService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.ZoneListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("zone service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("zone service: list: count: %w", err)
	}

	data := make([]dto.ZoneResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.ZoneListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *ZoneService) Update(ctx context.Context, id, orgID string, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("zone service: update: find: %w", err)
	}
	e.BuildingId = req.BuildingId
	e.Floor = req.Floor
	e.Name = req.Name

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("zone service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("zone service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ZoneService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("zone service: delete: %w", err)
	}
	return nil
}

func (s *ZoneService) toResponse(e *entity.Zone) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		BuildingId:     e.BuildingId,
		Floor:          e.Floor,
		Name:           e.Name,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
