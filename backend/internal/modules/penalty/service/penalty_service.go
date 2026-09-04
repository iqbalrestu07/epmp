package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/penalty/dto"
	"github.com/epmp/backend/internal/modules/penalty/entity"
	"github.com/epmp/backend/internal/modules/penalty/repository"
)

// PenaltyService implements the application layer for Penalty.
type PenaltyService struct {
	repo repository.PenaltyRepository
}

// NewPenaltyService creates a new PenaltyService.
func NewPenaltyService(repo repository.PenaltyRepository) *PenaltyService {
	return &PenaltyService{repo: repo}
}

func (s *PenaltyService) Create(ctx context.Context, orgID string, req *dto.CreatePenaltyRequest) (*dto.PenaltyResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("penalty service: create: organization ID required")
	}
	e := entity.NewPenalty()
	e.OrganizationId = orgID
	e.InvoiceId = req.InvoiceId
	e.Amount = req.Amount
	e.Status = req.Status
	e.PenaltyDate = req.PenaltyDate
	e.Description = req.Description

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("penalty service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *PenaltyService) GetByID(ctx context.Context, id, orgID string) (*dto.PenaltyResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("penalty service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *PenaltyService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.PenaltyListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("penalty service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("penalty service: list: count: %w", err)
	}

	data := make([]dto.PenaltyResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.PenaltyListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *PenaltyService) Update(ctx context.Context, id, orgID string, req *dto.UpdatePenaltyRequest) (*dto.PenaltyResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("penalty service: update: find: %w", err)
	}
	e.InvoiceId = req.InvoiceId
	e.Amount = req.Amount
	e.Status = req.Status
	e.PenaltyDate = req.PenaltyDate
	e.Description = req.Description

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("penalty service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("penalty service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *PenaltyService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("penalty service: delete: %w", err)
	}
	return nil
}

func (s *PenaltyService) toResponse(e *entity.Penalty) *dto.PenaltyResponse {
	return &dto.PenaltyResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		InvoiceId: e.InvoiceId,
		Amount: e.Amount,
		Status: e.Status,
		PenaltyDate: e.PenaltyDate,
		Description: e.Description,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
