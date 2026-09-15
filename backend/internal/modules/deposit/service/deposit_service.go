package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/deposit/dto"
	"github.com/epmp/backend/internal/modules/deposit/entity"
	"github.com/epmp/backend/internal/modules/deposit/repository"
)

// DepositService implements the application layer for Deposit.
type DepositService struct {
	repo repository.DepositRepository
}

// NewDepositService creates a new DepositService.
func NewDepositService(repo repository.DepositRepository) *DepositService {
	return &DepositService{repo: repo}
}

func (s *DepositService) Create(ctx context.Context, orgID string, req *dto.CreateDepositRequest) (*dto.DepositResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("deposit service: create: organization ID required")
	}
	e := entity.NewDeposit()
	e.OrganizationId = orgID
	e.ContractId = req.ContractId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.Status = req.Status
	e.CollectionDate = req.CollectionDate
	e.RefundDate = req.RefundDate
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("deposit service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *DepositService) GetByID(ctx context.Context, id, orgID string) (*dto.DepositResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("deposit service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *DepositService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.DepositListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("deposit service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("deposit service: list: count: %w", err)
	}

	data := make([]dto.DepositResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.DepositListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *DepositService) Update(ctx context.Context, id, orgID string, req *dto.UpdateDepositRequest) (*dto.DepositResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("deposit service: update: find: %w", err)
	}
	e.ContractId = req.ContractId
	e.TenantId = req.TenantId
	e.Amount = req.Amount
	e.Status = req.Status
	e.CollectionDate = req.CollectionDate
	e.RefundDate = req.RefundDate
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("deposit service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("deposit service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *DepositService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("deposit service: delete: %w", err)
	}
	return nil
}

func (s *DepositService) toResponse(e *entity.Deposit) *dto.DepositResponse {
	return &dto.DepositResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		ContractId:     e.ContractId,
		TenantId:       e.TenantId,
		Amount:         e.Amount,
		Status:         e.Status,
		CollectionDate: e.CollectionDate,
		RefundDate:     e.RefundDate,
		Notes:          e.Notes,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
