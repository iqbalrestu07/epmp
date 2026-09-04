package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/assetinspection/dto"
	"github.com/epmp/backend/internal/modules/assetinspection/entity"
	"github.com/epmp/backend/internal/modules/assetinspection/repository"
)

// AssetInspectionService implements the application layer for AssetInspection.
type AssetInspectionService struct {
	repo repository.AssetInspectionRepository
}

// NewAssetInspectionService creates a new AssetInspectionService.
func NewAssetInspectionService(repo repository.AssetInspectionRepository) *AssetInspectionService {
	return &AssetInspectionService{repo: repo}
}

func (s *AssetInspectionService) Create(ctx context.Context, orgID string, req *dto.CreateAssetInspectionRequest) (*dto.AssetInspectionResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("assetinspection service: create: organization ID required")
	}
	e := entity.NewAssetInspection()
	e.OrganizationId = orgID
	e.AssetId = req.AssetId
	e.InspectionDate = req.InspectionDate
	e.Condition = req.Condition
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("assetinspection service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetInspectionService) GetByID(ctx context.Context, id, orgID string) (*dto.AssetInspectionResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetinspection service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *AssetInspectionService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.AssetInspectionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetinspection service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetinspection service: list: count: %w", err)
	}

	data := make([]dto.AssetInspectionResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.AssetInspectionListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *AssetInspectionService) Update(ctx context.Context, id, orgID string, req *dto.UpdateAssetInspectionRequest) (*dto.AssetInspectionResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetinspection service: update: find: %w", err)
	}
	e.AssetId = req.AssetId
	e.InspectionDate = req.InspectionDate
	e.Condition = req.Condition
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("assetinspection service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("assetinspection service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetInspectionService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("assetinspection service: delete: %w", err)
	}
	return nil
}

func (s *AssetInspectionService) toResponse(e *entity.AssetInspection) *dto.AssetInspectionResponse {
	return &dto.AssetInspectionResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		AssetId: e.AssetId,
		InspectionDate: e.InspectionDate,
		Condition: e.Condition,
		Notes: e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
