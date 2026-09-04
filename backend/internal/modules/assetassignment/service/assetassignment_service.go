package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/assetassignment/dto"
	"github.com/epmp/backend/internal/modules/assetassignment/entity"
	"github.com/epmp/backend/internal/modules/assetassignment/repository"
)

// AssetAssignmentService implements the application layer for AssetAssignment.
type AssetAssignmentService struct {
	repo repository.AssetAssignmentRepository
}

// NewAssetAssignmentService creates a new AssetAssignmentService.
func NewAssetAssignmentService(repo repository.AssetAssignmentRepository) *AssetAssignmentService {
	return &AssetAssignmentService{repo: repo}
}

func (s *AssetAssignmentService) Create(ctx context.Context, orgID string, req *dto.CreateAssetAssignmentRequest) (*dto.AssetAssignmentResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("assetassignment service: create: organization ID required")
	}
	e := entity.NewAssetAssignment()
	e.OrganizationId = orgID
	e.AssetId = req.AssetId
	e.RoomId = req.RoomId
	e.AssignedDate = req.AssignedDate

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("assetassignment service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetAssignmentService) GetByID(ctx context.Context, id, orgID string) (*dto.AssetAssignmentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetassignment service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *AssetAssignmentService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.AssetAssignmentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetassignment service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetassignment service: list: count: %w", err)
	}

	data := make([]dto.AssetAssignmentResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.AssetAssignmentListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *AssetAssignmentService) Update(ctx context.Context, id, orgID string, req *dto.UpdateAssetAssignmentRequest) (*dto.AssetAssignmentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetassignment service: update: find: %w", err)
	}
	e.AssetId = req.AssetId
	e.RoomId = req.RoomId
	e.AssignedDate = req.AssignedDate

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("assetassignment service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetassignment service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetAssignmentService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("assetassignment service: delete: %w", err)
	}
	return nil
}

func (s *AssetAssignmentService) toResponse(e *entity.AssetAssignment) *dto.AssetAssignmentResponse {
	return &dto.AssetAssignmentResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		AssetId: e.AssetId,
		RoomId: e.RoomId,
		AssignedDate: e.AssignedDate,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
