package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/tenantcontact/dto"
	"github.com/epmp/backend/internal/modules/tenantcontact/entity"
	"github.com/epmp/backend/internal/modules/tenantcontact/repository"
)

// TenantContactService implements the application layer for TenantContact.
type TenantContactService struct {
	repo repository.TenantContactRepository
}

// NewTenantContactService creates a new TenantContactService.
func NewTenantContactService(repo repository.TenantContactRepository) *TenantContactService {
	return &TenantContactService{repo: repo}
}

func (s *TenantContactService) Create(ctx context.Context, orgID string, req *dto.CreateTenantContactRequest) (*dto.TenantContactResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("tenantcontact service: create: organization ID required")
	}
	e := entity.NewTenantContact()
	e.OrganizationId = orgID
	e.TenantId = req.TenantId
	e.ContactType = req.ContactType
	e.ContactValue = req.ContactValue
	e.IsPrimary = req.IsPrimary

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantcontact service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantContactService) GetByID(ctx context.Context, id, orgID string) (*dto.TenantContactResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *TenantContactService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.TenantContactListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact service: list: count: %w", err)
	}

	data := make([]dto.TenantContactResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.TenantContactListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TenantContactService) Update(ctx context.Context, id, orgID string, req *dto.UpdateTenantContactRequest) (*dto.TenantContactResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact service: update: find: %w", err)
	}
	e.TenantId = req.TenantId
	e.ContactType = req.ContactType
	e.ContactValue = req.ContactValue
	e.IsPrimary = req.IsPrimary

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantcontact service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantContactService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("tenantcontact service: delete: %w", err)
	}
	return nil
}

func (s *TenantContactService) toResponse(e *entity.TenantContact) *dto.TenantContactResponse {
	return &dto.TenantContactResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		TenantId: e.TenantId,
		ContactType: e.ContactType,
		ContactValue: e.ContactValue,
		IsPrimary: e.IsPrimary,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
