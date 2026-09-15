package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/contract/dto"
	"github.com/epmp/backend/internal/modules/contract/entity"
	"github.com/epmp/backend/internal/modules/contract/repository"
)

// ContractService implements the application layer for Contract.
type ContractService struct {
	repo repository.ContractRepository
}

// NewContractService creates a new ContractService.
func NewContractService(repo repository.ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) Create(ctx context.Context, orgID string, req *dto.CreateContractRequest) (*dto.ContractResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("contract service: create: organization ID required")
	}
	e := entity.NewContract()
	e.OrganizationId = orgID
	e.ReservationId = req.ReservationId
	e.TenantId = req.TenantId
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Status = req.Status
	e.StartDate = req.StartDate
	e.EndDate = req.EndDate
	e.MonthlyRent = req.MonthlyRent
	e.DepositAmount = req.DepositAmount
	e.Terms = req.Terms

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("contract service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ContractService) GetByID(ctx context.Context, id, orgID string) (*dto.ContractResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("contract service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *ContractService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.ContractListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("contract service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("contract service: list: count: %w", err)
	}

	data := make([]dto.ContractResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.ContractListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *ContractService) Update(ctx context.Context, id, orgID string, req *dto.UpdateContractRequest) (*dto.ContractResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("contract service: update: find: %w", err)
	}
	e.ReservationId = req.ReservationId
	e.TenantId = req.TenantId
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Status = req.Status
	e.StartDate = req.StartDate
	e.EndDate = req.EndDate
	e.MonthlyRent = req.MonthlyRent
	e.DepositAmount = req.DepositAmount
	e.Terms = req.Terms

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("contract service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("contract service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ContractService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("contract service: delete: %w", err)
	}
	return nil
}

func (s *ContractService) toResponse(e *entity.Contract) *dto.ContractResponse {
	return &dto.ContractResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		ReservationId:  e.ReservationId,
		TenantId:       e.TenantId,
		PropertyId:     e.PropertyId,
		RoomId:         e.RoomId,
		Status:         e.Status,
		StartDate:      e.StartDate,
		EndDate:        e.EndDate,
		MonthlyRent:    e.MonthlyRent,
		DepositAmount:  e.DepositAmount,
		Terms:          e.Terms,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
