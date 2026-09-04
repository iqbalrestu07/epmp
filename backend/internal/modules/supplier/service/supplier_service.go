package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/supplier/dto"
	"github.com/epmp/backend/internal/modules/supplier/entity"
	"github.com/epmp/backend/internal/modules/supplier/repository"
)

// SupplierService implements the application layer for Supplier.
type SupplierService struct {
	repo repository.SupplierRepository
}

// NewSupplierService creates a new SupplierService.
func NewSupplierService(repo repository.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) Create(ctx context.Context, orgID string, req *dto.CreateSupplierRequest) (*dto.SupplierResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("supplier service: create: organization ID required")
	}
	e := entity.NewSupplier()
	e.OrganizationId = orgID
	e.Name = req.Name
	e.ContactPerson = req.ContactPerson
	e.Phone = req.Phone
	e.ServiceType = req.ServiceType

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("supplier service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *SupplierService) GetByID(ctx context.Context, id, orgID string) (*dto.SupplierResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("supplier service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *SupplierService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.SupplierListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("supplier service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("supplier service: list: count: %w", err)
	}

	data := make([]dto.SupplierResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.SupplierListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *SupplierService) Update(ctx context.Context, id, orgID string, req *dto.UpdateSupplierRequest) (*dto.SupplierResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("supplier service: update: find: %w", err)
	}
	e.Name = req.Name
	e.ContactPerson = req.ContactPerson
	e.Phone = req.Phone
	e.ServiceType = req.ServiceType

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("supplier service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("supplier service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *SupplierService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("supplier service: delete: %w", err)
	}
	return nil
}

func (s *SupplierService) toResponse(e *entity.Supplier) *dto.SupplierResponse {
	return &dto.SupplierResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		Name: e.Name,
		ContactPerson: e.ContactPerson,
		Phone: e.Phone,
		ServiceType: e.ServiceType,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
