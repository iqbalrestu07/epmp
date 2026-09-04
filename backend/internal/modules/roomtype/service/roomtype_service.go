package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/roomtype/dto"
	"github.com/epmp/backend/internal/modules/roomtype/entity"
	"github.com/epmp/backend/internal/modules/roomtype/repository"
)

// RoomTypeService implements the application layer for RoomType.
type RoomTypeService struct {
	repo repository.RoomTypeRepository
}

// NewRoomTypeService creates a new RoomTypeService.
func NewRoomTypeService(repo repository.RoomTypeRepository) *RoomTypeService {
	return &RoomTypeService{repo: repo}
}

func (s *RoomTypeService) Create(ctx context.Context, orgID string, req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("roomtype service: create: organization ID required")
	}
	e := entity.NewRoomType()
	e.OrganizationId = orgID
	e.Name = req.Name
	e.Description = req.Description
	e.BasePrice = req.BasePrice

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("roomtype service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *RoomTypeService) GetByID(ctx context.Context, id, orgID string) (*dto.RoomTypeResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("roomtype service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *RoomTypeService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.RoomTypeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("roomtype service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("roomtype service: list: count: %w", err)
	}

	data := make([]dto.RoomTypeResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.RoomTypeListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *RoomTypeService) Update(ctx context.Context, id, orgID string, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("roomtype service: update: find: %w", err)
	}
	e.Name = req.Name
	e.Description = req.Description
	e.BasePrice = req.BasePrice

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("roomtype service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("roomtype service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *RoomTypeService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("roomtype service: delete: %w", err)
	}
	return nil
}

func (s *RoomTypeService) toResponse(e *entity.RoomType) *dto.RoomTypeResponse {
	return &dto.RoomTypeResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		Name: e.Name,
		Description: e.Description,
		BasePrice: e.BasePrice,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
