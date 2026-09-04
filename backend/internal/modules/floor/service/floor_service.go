package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/floor/dto"
	"github.com/epmp/backend/internal/modules/floor/entity"
	"github.com/epmp/backend/internal/modules/floor/repository"
)

// FloorService implements the application layer for Floor.
type FloorService struct {
	repo repository.FloorRepository
}

// NewFloorService creates a new FloorService.
func NewFloorService(repo repository.FloorRepository) *FloorService {
	return &FloorService{repo: repo}
}

func (s *FloorService) Create(ctx context.Context, req *dto.CreateFloorRequest) (*dto.FloorResponse, error) {
	e := entity.NewFloor()
	e.OrganizationId = req.OrganizationId
	e.BuildingId = req.BuildingId
	e.Name = req.Name
	e.FloorNumber = req.FloorNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("floor service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *FloorService) GetByID(ctx context.Context, id string) (*dto.FloorResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("floor service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *FloorService) List(ctx context.Context, page, perPage int, search string, buildingId string) (*dto.FloorListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, buildingId)
	if err != nil {
		return nil, fmt.Errorf("floor service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, buildingId)
	if err != nil {
		return nil, fmt.Errorf("floor service: list: count: %w", err)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	data := make([]dto.FloorResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	return &dto.FloorListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *FloorService) Update(ctx context.Context, id string, req *dto.UpdateFloorRequest) (*dto.FloorResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("floor service: update: find: %w", err)
	}
	e.OrganizationId = req.OrganizationId
	e.BuildingId = req.BuildingId
	e.Name = req.Name
	e.FloorNumber = req.FloorNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("floor service: update: save: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.toResponse(e), nil
	}
	return s.toResponse(updated), nil
}

func (s *FloorService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("floor service: delete: %w", err)
	}
	return nil
}

func (s *FloorService) toResponse(e *entity.Floor) *dto.FloorResponse {
	return &dto.FloorResponse{
		Id:             e.Id,
		OrganizationId: e.OrganizationId,
		BuildingId:     e.BuildingId,
		Name:           e.Name,
		FloorNumber:    e.FloorNumber,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
