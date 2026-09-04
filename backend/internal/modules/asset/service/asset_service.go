package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/asset/dto"
	"github.com/epmp/backend/internal/modules/asset/entity"
	"github.com/epmp/backend/internal/modules/asset/repository"
	"github.com/oklog/ulid/v2"
)

// AssetService implements the application layer for Asset.
type AssetService struct {
	repo repository.AssetRepository
}

// NewAssetService creates a new AssetService.
func NewAssetService(repo repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) Create(ctx context.Context, orgID string, req *dto.CreateAssetRequest) (*dto.AssetResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("asset service: create: organization ID required")
	}
	e := entity.NewAsset()
	e.Id = ulid.Make().String()
	e.OrganizationId = orgID
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.Category = req.Category
	e.Status = req.Status
	e.PurchasePrice = req.PurchasePrice

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("asset service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetService) GetByID(ctx context.Context, id, orgID string) (*dto.AssetResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("asset service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *AssetService) List(ctx context.Context, page, perPage int, search, propertyID, orgID string) (*dto.AssetListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, propertyID, orgID)
	if err != nil {
		return nil, fmt.Errorf("asset service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, propertyID, orgID)
	if err != nil {
		return nil, fmt.Errorf("asset service: list: count: %w", err)
	}

	data := make([]dto.AssetResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.AssetListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *AssetService) Update(ctx context.Context, id, orgID string, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("asset service: update: find: %w", err)
	}
	e.PropertyId = req.PropertyId
	e.Name = req.Name
	e.Category = req.Category
	e.Status = req.Status
	e.PurchasePrice = req.PurchasePrice

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("asset service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("asset service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *AssetService) Delete(ctx context.Context, id, orgID string) error {
	// First ensure it exists and belongs to the org
	_, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return fmt.Errorf("asset service: delete: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("asset service: delete: %w", err)
	}
	return nil
}

func (s *AssetService) toResponse(e *entity.Asset) *dto.AssetResponse {
	return &dto.AssetResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		PropertyId:     e.PropertyId,
		Name:           e.Name,
		Category:       e.Category,
		Status:         e.Status,
		PurchasePrice:  e.PurchasePrice,
	}
}
