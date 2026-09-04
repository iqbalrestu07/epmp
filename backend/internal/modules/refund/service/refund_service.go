package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/refund/dto"
	"github.com/epmp/backend/internal/modules/refund/entity"
	"github.com/epmp/backend/internal/modules/refund/repository"
)

// RefundService implements the application layer for Refund.
type RefundService struct {
	repo repository.RefundRepository
}

// NewRefundService creates a new RefundService.
func NewRefundService(repo repository.RefundRepository) *RefundService {
	return &RefundService{repo: repo}
}

func (s *RefundService) Create(ctx context.Context, orgID string, req *dto.CreateRefundRequest) (*dto.RefundResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("refund service: create: organization ID required")
	}
	e := entity.NewRefund()
	e.OrganizationId = orgID
	e.PaymentId = req.PaymentId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.Status = req.Status
	e.RefundDate = req.RefundDate
	e.Reason = req.Reason

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("refund service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *RefundService) GetByID(ctx context.Context, id, orgID string) (*dto.RefundResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("refund service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *RefundService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.RefundListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("refund service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("refund service: list: count: %w", err)
	}

	data := make([]dto.RefundResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.RefundListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *RefundService) Update(ctx context.Context, id, orgID string, req *dto.UpdateRefundRequest) (*dto.RefundResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("refund service: update: find: %w", err)
	}
	e.PaymentId = req.PaymentId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.Status = req.Status
	e.RefundDate = req.RefundDate
	e.Reason = req.Reason

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("refund service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("refund service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *RefundService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("refund service: delete: %w", err)
	}
	return nil
}

func (s *RefundService) toResponse(e *entity.Refund) *dto.RefundResponse {
	return &dto.RefundResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		PaymentId: e.PaymentId,
		TenantId: e.TenantId,
		Amount: e.Amount,
		Status: e.Status,
		RefundDate: e.RefundDate,
		Reason: e.Reason,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
