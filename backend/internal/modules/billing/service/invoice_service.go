package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/billing/dto"
	"github.com/epmp/backend/internal/modules/billing/entity"
	"github.com/epmp/backend/internal/modules/billing/repository"
)

// InvoiceService implements the application layer for Invoice.
type InvoiceService struct {
	repo repository.InvoiceRepository
}

// NewInvoiceService creates a new InvoiceService.
func NewInvoiceService(repo repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) Create(ctx context.Context, orgID string, req *dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("billing service: create: organization ID required")
	}
	e := entity.NewInvoice()
	e.OrganizationId = orgID
	e.ContractId = req.ContractId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.Currency = req.Currency
	e.Status = req.Status
	e.DueDate = req.DueDate
	e.PaidDate = req.PaidDate
	e.PaymentMethod = req.PaymentMethod
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("billing service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *InvoiceService) GetByID(ctx context.Context, id, orgID string) (*dto.InvoiceResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("billing service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *InvoiceService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("billing service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("billing service: list: count: %w", err)
	}

	data := make([]dto.InvoiceResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.InvoiceListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *InvoiceService) Update(ctx context.Context, id, orgID string, req *dto.UpdateInvoiceRequest) (*dto.InvoiceResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("billing service: update: find: %w", err)
	}
	e.ContractId = req.ContractId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	if req.Currency != "" {
		e.Currency = req.Currency
	}
	e.Status = req.Status
	e.DueDate = req.DueDate
	e.PaidDate = req.PaidDate
	e.PaymentMethod = req.PaymentMethod
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("billing service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("billing service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *InvoiceService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("billing service: delete: %w", err)
	}
	return nil
}

func (s *InvoiceService) toResponse(e *entity.Invoice) *dto.InvoiceResponse {
	return &dto.InvoiceResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		ContractId:     e.ContractId,
		TenantId:       e.TenantId,
		Amount:         e.Amount,
		Currency:       e.Currency,
		Status:         e.Status,
		DueDate:        e.DueDate,
		PaidDate:       e.PaidDate,
		PaymentMethod:  e.PaymentMethod,
		Notes:          e.Notes,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
