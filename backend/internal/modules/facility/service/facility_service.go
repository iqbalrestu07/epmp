package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/facility/dto"
	"github.com/epmp/backend/internal/modules/facility/entity"
	"github.com/epmp/backend/internal/modules/facility/repository"
)

// FacilityService implements the application layer for Facility.
type FacilityService struct {
	repo repository.FacilityRepository
}

// NewFacilityService creates a new FacilityService.
func NewFacilityService(repo repository.FacilityRepository) *FacilityService {
	return &FacilityService{repo: repo}
}

func (s *FacilityService) Create(ctx context.Context, orgID string, req *dto.CreateFacilityRequest) (*dto.FacilityResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("facility service: create: organization ID required")
	}
	e := entity.NewFacility()
	e.OrganizationId = orgID
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.Description = req.Description

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("facility service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *FacilityService) GetByID(ctx context.Context, id, orgID string) (*dto.FacilityResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("facility service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *FacilityService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.FacilityListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("facility service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("facility service: list: count: %w", err)
	}

	data := make([]dto.FacilityResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.FacilityListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *FacilityService) Update(ctx context.Context, id, orgID string, req *dto.UpdateFacilityRequest) (*dto.FacilityResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("facility service: update: find: %w", err)
	}
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.Description = req.Description

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("facility service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("facility service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *FacilityService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("facility service: delete: %w", err)
	}
	return nil
}

func (s *FacilityService) toResponse(e *entity.Facility) *dto.FacilityResponse {
	return &dto.FacilityResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		PropertyId:     e.PropertyId,
		Name:           e.Name,
		Description:    e.Description,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
