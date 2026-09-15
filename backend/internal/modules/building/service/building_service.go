package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/building/dto"
	"github.com/epmp/backend/internal/modules/building/entity"
	"github.com/epmp/backend/internal/modules/building/repository"
)

// BuildingService implements the application layer for Building.
type BuildingService struct {
	repo repository.BuildingRepository
}

// NewBuildingService creates a new BuildingService.
func NewBuildingService(repo repository.BuildingRepository) *BuildingService {
	return &BuildingService{repo: repo}
}

func (s *BuildingService) Create(ctx context.Context, req *dto.CreateBuildingRequest, orgID string) (*dto.BuildingResponse, error) {
	e := entity.NewBuilding()
	e.OrganizationId = orgID
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.TotalFloors = req.TotalFloors

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("building service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *BuildingService) GetByID(ctx context.Context, id, orgID string) (*dto.BuildingResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("building service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *BuildingService) List(ctx context.Context, page, perPage int, search string, propertyId, orgID string) (*dto.BuildingListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, propertyId, orgID)
	if err != nil {
		return nil, fmt.Errorf("building service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, propertyId, orgID)
	if err != nil {
		return nil, fmt.Errorf("building service: list: count: %w", err)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	data := make([]dto.BuildingResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	return &dto.BuildingListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *BuildingService) Update(ctx context.Context, id, orgID string, req *dto.UpdateBuildingRequest) (*dto.BuildingResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("building service: update: find: %w", err)
	}
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.TotalFloors = req.TotalFloors

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("building service: update: save: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return s.toResponse(e), nil
	}
	return s.toResponse(updated), nil
}

func (s *BuildingService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("building service: delete: %w", err)
	}
	return nil
}

func (s *BuildingService) toResponse(e *entity.Building) *dto.BuildingResponse {
	return &dto.BuildingResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		PropertyId:     e.PropertyId,
		Name:           e.Name,
		TotalFloors:    e.TotalFloors,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
