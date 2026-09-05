package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/room/dto"
	"github.com/epmp/backend/internal/modules/room/entity"
	"github.com/epmp/backend/internal/modules/room/repository"
	"github.com/epmp/backend/internal/pkg/uid"
)

// RoomService implements the application layer for Room.
type RoomService struct {
	repo repository.RoomRepository
}

// NewRoomService creates a new RoomService.
func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Create(ctx context.Context, req *dto.CreateRoomRequest, orgID string) (*dto.RoomResponse, error) {
	e := entity.NewRoom()
	e.Id = uid.New()
	e.OrganizationId = orgID
	e.PropertyId = req.PropertyId
	e.FloorId = req.FloorId
	e.RoomTypeId = req.RoomTypeId
	e.Name = req.Name
	e.Capacity = req.Capacity
	e.Price = req.Price
	e.IsAvailable = req.IsAvailable

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("room service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *RoomService) GetByID(ctx context.Context, id, orgID string) (*dto.RoomResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("room service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *RoomService) List(ctx context.Context, page, perPage int, search, floorId, propertyId, buildingId, orgID string) (*dto.RoomListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, floorId, propertyId, buildingId, orgID)
	if err != nil {
		return nil, fmt.Errorf("room service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, floorId, propertyId, buildingId, orgID)
	if err != nil {
		return nil, fmt.Errorf("room service: list: count: %w", err)
	}

	data := make([]dto.RoomResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return &dto.RoomListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *RoomService) Update(ctx context.Context, id, orgID string, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("room service: update: find: %w", err)
	}
	e.PropertyId = req.PropertyId
	e.FloorId = req.FloorId
	e.RoomTypeId = req.RoomTypeId
	e.Name = req.Name
	e.Capacity = req.Capacity
	e.Price = req.Price
	e.IsAvailable = req.IsAvailable

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("room service: update: save: %w", err)
	}

	updated, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return s.toResponse(e), nil
	}
	return s.toResponse(updated), nil
}

func (s *RoomService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("room service: delete: %w", err)
	}
	return nil
}

func (s *RoomService) toResponse(e *entity.Room) *dto.RoomResponse {
	return &dto.RoomResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		PropertyId:     e.PropertyId,
		FloorId:        e.FloorId,
		RoomTypeId:     e.RoomTypeId,
		Name:           e.Name,
		Capacity:       e.Capacity,
		Price:          e.Price,
		IsAvailable:    e.IsAvailable,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
