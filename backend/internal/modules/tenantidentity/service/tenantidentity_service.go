package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/tenantidentity/dto"
	"github.com/epmp/backend/internal/modules/tenantidentity/entity"
	"github.com/epmp/backend/internal/modules/tenantidentity/repository"
)

// TenantIdentityService implements the application layer for TenantIdentity.
type TenantIdentityService struct {
	repo repository.TenantIdentityRepository
}

// NewTenantIdentityService creates a new TenantIdentityService.
func NewTenantIdentityService(repo repository.TenantIdentityRepository) *TenantIdentityService {
	return &TenantIdentityService{repo: repo}
}

func (s *TenantIdentityService) Create(ctx context.Context, orgID string, req *dto.CreateTenantIdentityRequest) (*dto.TenantIdentityResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("tenantidentity service: create: organization ID required")
	}
	e := entity.NewTenantIdentity()
	e.OrganizationId = orgID
	e.TenantId = req.TenantId
	e.IdentityType = req.IdentityType
	e.IdentityNumber = req.IdentityNumber
	e.FileUrl = req.FileUrl

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantidentity service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantIdentityService) GetByID(ctx context.Context, id, orgID string) (*dto.TenantIdentityResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *TenantIdentityService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.TenantIdentityListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity service: list: count: %w", err)
	}

	data := make([]dto.TenantIdentityResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.TenantIdentityListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TenantIdentityService) Update(ctx context.Context, id, orgID string, req *dto.UpdateTenantIdentityRequest) (*dto.TenantIdentityResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity service: update: find: %w", err)
	}
	e.TenantId = req.TenantId
	e.IdentityType = req.IdentityType
	e.IdentityNumber = req.IdentityNumber
	e.FileUrl = req.FileUrl

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantidentity service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantIdentityService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("tenantidentity service: delete: %w", err)
	}
	return nil
}

func (s *TenantIdentityService) toResponse(e *entity.TenantIdentity) *dto.TenantIdentityResponse {
	return &dto.TenantIdentityResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		TenantId: e.TenantId,
		IdentityType: e.IdentityType,
		IdentityNumber: e.IdentityNumber,
		FileUrl: e.FileUrl,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
