package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/tenant/dto"
	"github.com/epmp/backend/internal/modules/tenant/entity"
	"github.com/epmp/backend/internal/modules/tenant/repository"
)

// TenantService implements the application layer for Tenant.
type TenantService struct {
	repo repository.TenantRepository
}

// NewTenantService creates a new TenantService.
func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) Create(ctx context.Context, orgID string, req *dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("tenant service: create: organization ID required")
	}
	e := entity.NewTenant()
	e.OrganizationId = orgID
	e.FullName = req.FullName
	e.Email = req.Email
	e.Phone = req.Phone
	e.IdentityNumber = req.IdentityNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenant service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantService) GetByID(ctx context.Context, id, orgID string) (*dto.TenantResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenant service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *TenantService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.TenantListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenant service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenant service: list: count: %w", err)
	}

	data := make([]dto.TenantResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.TenantListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TenantService) Update(ctx context.Context, id, orgID string, req *dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenant service: update: find: %w", err)
	}
	e.FullName = req.FullName
	e.Email = req.Email
	e.Phone = req.Phone
	e.IdentityNumber = req.IdentityNumber
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenant service: update: save: %w", err)
	}

	// Re-fetch to get updated timestamps
	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenant service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("tenant service: delete: %w", err)
	}
	return nil
}

func (s *TenantService) toResponse(e *entity.Tenant) *dto.TenantResponse {
	return &dto.TenantResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		FullName:       e.FullName,
		Email:          e.Email,
		Phone:          e.Phone,
		IdentityNumber: e.IdentityNumber,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
