package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/payment/dto"
	"github.com/epmp/backend/internal/modules/payment/entity"
	"github.com/epmp/backend/internal/modules/payment/repository"
)

// PaymentService implements the application layer for Payment.
type PaymentService struct {
	repo repository.PaymentRepository
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(ctx context.Context, orgID string, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("payment service: create: organization ID required")
	}
	e := entity.NewPayment()
	e.OrganizationId = orgID
	e.InvoiceId = req.InvoiceId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.PaymentDate = req.PaymentDate
	e.PaymentMethod = req.PaymentMethod
	e.Status = req.Status
	e.ReferenceNumber = req.ReferenceNumber

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("payment service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *PaymentService) GetByID(ctx context.Context, id, orgID string) (*dto.PaymentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("payment service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *PaymentService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.PaymentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("payment service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("payment service: list: count: %w", err)
	}

	data := make([]dto.PaymentResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.PaymentListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *PaymentService) Update(ctx context.Context, id, orgID string, req *dto.UpdatePaymentRequest) (*dto.PaymentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("payment service: update: find: %w", err)
	}
	e.InvoiceId = req.InvoiceId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.PaymentDate = req.PaymentDate
	e.PaymentMethod = req.PaymentMethod
	e.Status = req.Status
	e.ReferenceNumber = req.ReferenceNumber

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("payment service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("payment service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *PaymentService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("payment service: delete: %w", err)
	}
	return nil
}

func (s *PaymentService) toResponse(e *entity.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		InvoiceId: e.InvoiceId,
		TenantId: e.TenantId,
		Amount: e.Amount,
		PaymentDate: e.PaymentDate,
		PaymentMethod: e.PaymentMethod,
		Status: e.Status,
		ReferenceNumber: e.ReferenceNumber,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
