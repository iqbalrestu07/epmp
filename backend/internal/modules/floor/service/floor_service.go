package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/building/repository"
	"github.com/epmp/backend/internal/modules/floor/dto"
	"github.com/epmp/backend/internal/modules/floor/entity"
	floorrepo "github.com/epmp/backend/internal/modules/floor/repository"
)

// FloorService implements the application layer for Floor.
type FloorService struct {
	repo         floorrepo.FloorRepository
	buildingRepo repository.BuildingRepository
}

// NewFloorService creates a new FloorService.
func NewFloorService(repo floorrepo.FloorRepository, buildingRepo repository.BuildingRepository) *FloorService {
	return &FloorService{repo: repo, buildingRepo: buildingRepo}
}

func (s *FloorService) Create(ctx context.Context, req *dto.CreateFloorRequest, orgID string) (*dto.FloorResponse, error) {
	// Validate against building's total_floors limit
	if s.buildingRepo != nil && req.BuildingId != "" {
		building, err := s.buildingRepo.FindByID(ctx, req.BuildingId, orgID)
		if err != nil {
			return nil, fmt.Errorf("floor service: create: building not found: %w", err)
		}
		count, err := s.repo.CountByBuildingID(ctx, req.BuildingId)
		if err != nil {
			return nil, fmt.Errorf("floor service: create: count floors: %w", err)
		}
		if int(count) >= building.TotalFloors {
			return nil, fmt.Errorf("floor service: create: cannot add more floors than building's total_floors (%d)", building.TotalFloors)
		}
	}

	e := entity.NewFloor()
	e.OrganizationId = orgID
	e.BuildingId = req.BuildingId
	e.Name = req.Name
	e.FloorNumber = req.FloorNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("floor service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *FloorService) GetByID(ctx context.Context, id, orgID string) (*dto.FloorResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("floor service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *FloorService) List(ctx context.Context, page, perPage int, search string, buildingId, orgID string) (*dto.FloorListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, buildingId, orgID)
	if err != nil {
		return nil, fmt.Errorf("floor service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, buildingId, orgID)
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

func (s *FloorService) Update(ctx context.Context, id, orgID string, req *dto.UpdateFloorRequest) (*dto.FloorResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("floor service: update: find: %w", err)
	}
	e.BuildingId = req.BuildingId
	e.Name = req.Name
	e.FloorNumber = req.FloorNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("floor service: update: save: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return s.toResponse(e), nil
	}
	return s.toResponse(updated), nil
}

func (s *FloorService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
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
