package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/charge/dto"
	"github.com/epmp/backend/internal/modules/charge/entity"
	"github.com/epmp/backend/internal/modules/charge/repository"
)

// ChargeService implements the application layer for Charge.
type ChargeService struct {
	repo repository.ChargeRepository
}

// NewChargeService creates a new ChargeService.
func NewChargeService(repo repository.ChargeRepository) *ChargeService {
	return &ChargeService{repo: repo}
}

func (s *ChargeService) Create(ctx context.Context, orgID string, req *dto.CreateChargeRequest) (*dto.ChargeResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("charge service: create: organization ID required")
	}
	e := entity.NewCharge()
	e.OrganizationId = orgID
	e.ContractId = req.ContractId
	e.InvoiceId = req.InvoiceId
	e.ChargeType = req.ChargeType
	e.Amount = req.Amount
	e.Status = req.Status
	e.ChargeDate = req.ChargeDate
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("charge service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ChargeService) GetByID(ctx context.Context, id, orgID string) (*dto.ChargeResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("charge service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *ChargeService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.ChargeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("charge service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("charge service: list: count: %w", err)
	}

	data := make([]dto.ChargeResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.ChargeListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *ChargeService) Update(ctx context.Context, id, orgID string, req *dto.UpdateChargeRequest) (*dto.ChargeResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("charge service: update: find: %w", err)
	}
	e.ContractId = req.ContractId
	e.InvoiceId = req.InvoiceId
	e.ChargeType = req.ChargeType
	e.Amount = req.Amount
	e.Status = req.Status
	e.ChargeDate = req.ChargeDate
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("charge service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("charge service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ChargeService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("charge service: delete: %w", err)
	}
	return nil
}

func (s *ChargeService) toResponse(e *entity.Charge) *dto.ChargeResponse {
	return &dto.ChargeResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		ContractId:     e.ContractId,
		InvoiceId:      e.InvoiceId,
		ChargeType:     e.ChargeType,
		Amount:         e.Amount,
		Status:         e.Status,
		ChargeDate:     e.ChargeDate,
		Notes:          e.Notes,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
