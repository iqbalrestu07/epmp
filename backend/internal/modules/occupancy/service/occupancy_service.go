package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/occupancy/dto"
	"github.com/epmp/backend/internal/modules/occupancy/entity"
	"github.com/epmp/backend/internal/modules/occupancy/repository"
)

// OccupancyService implements the application layer for Occupancy.
type OccupancyService struct {
	repo repository.OccupancyRepository
}

// NewOccupancyService creates a new OccupancyService.
func NewOccupancyService(repo repository.OccupancyRepository) *OccupancyService {
	return &OccupancyService{repo: repo}
}

func (s *OccupancyService) Create(ctx context.Context, orgID string, req *dto.CreateOccupancyRequest) (*dto.OccupancyResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("occupancy service: create: organization ID required")
	}
	e := entity.NewOccupancy()
	e.OrganizationId = orgID
	e.ContractId = req.ContractId
	e.RoomId = req.RoomId
	e.TenantId = req.TenantId
	e.Status = req.Status
	e.CheckInTime = req.CheckInTime
	e.CheckOutTime = req.CheckOutTime
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("occupancy service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *OccupancyService) GetByID(ctx context.Context, id, orgID string) (*dto.OccupancyResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("occupancy service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *OccupancyService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.OccupancyListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("occupancy service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("occupancy service: list: count: %w", err)
	}

	data := make([]dto.OccupancyResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.OccupancyListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *OccupancyService) Update(ctx context.Context, id, orgID string, req *dto.UpdateOccupancyRequest) (*dto.OccupancyResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("occupancy service: update: find: %w", err)
	}
	e.ContractId = req.ContractId
	e.RoomId = req.RoomId
	e.TenantId = req.TenantId
	e.Status = req.Status
	e.CheckInTime = req.CheckInTime
	e.CheckOutTime = req.CheckOutTime
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("occupancy service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("occupancy service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *OccupancyService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("occupancy service: delete: %w", err)
	}
	return nil
}

func (s *OccupancyService) toResponse(e *entity.Occupancy) *dto.OccupancyResponse {
	return &dto.OccupancyResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		ContractId: e.ContractId,
		RoomId: e.RoomId,
		TenantId: e.TenantId,
		Status: e.Status,
		CheckInTime: e.CheckInTime,
		CheckOutTime: e.CheckOutTime,
		Notes: e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
