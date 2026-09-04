package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/bed/dto"
	"github.com/epmp/backend/internal/modules/bed/entity"
	"github.com/epmp/backend/internal/modules/bed/repository"
)

// BedService implements the application layer for Bed.
type BedService struct {
	repo repository.BedRepository
}

// NewBedService creates a new BedService.
func NewBedService(repo repository.BedRepository) *BedService {
	return &BedService{repo: repo}
}

func (s *BedService) Create(ctx context.Context, orgID string, req *dto.CreateBedRequest) (*dto.BedResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("bed service: create: organization ID required")
	}
	e := entity.NewBed()
	e.OrganizationId = orgID
	e.RoomId = req.RoomId
	e.Name = req.Name
	e.Status = req.Status

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("bed service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *BedService) GetByID(ctx context.Context, id, orgID string) (*dto.BedResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("bed service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *BedService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.BedListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("bed service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("bed service: list: count: %w", err)
	}

	data := make([]dto.BedResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.BedListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *BedService) Update(ctx context.Context, id, orgID string, req *dto.UpdateBedRequest) (*dto.BedResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("bed service: update: find: %w", err)
	}
	e.RoomId = req.RoomId
	e.Name = req.Name
	e.Status = req.Status

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("bed service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("bed service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *BedService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("bed service: delete: %w", err)
	}
	return nil
}

func (s *BedService) toResponse(e *entity.Bed) *dto.BedResponse {
	return &dto.BedResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		RoomId: e.RoomId,
		Name: e.Name,
		Status: e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
