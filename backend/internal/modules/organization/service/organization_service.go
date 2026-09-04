package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/organization/dto"
	"github.com/epmp/backend/internal/modules/organization/entity"
	"github.com/epmp/backend/internal/modules/organization/repository"
)

// OrganizationService implements the application layer for Organization.
type OrganizationService struct {
	repo repository.OrganizationRepository
}

// NewOrganizationService creates a new OrganizationService.
func NewOrganizationService(repo repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) Create(ctx context.Context, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error) {
	e := entity.NewOrganization()
	e.Name = req.Name
	e.Domain = req.Domain
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("organization service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id string) (*dto.OrganizationResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("organization service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *OrganizationService) List(ctx context.Context, page, perPage int, search string) (*dto.OrganizationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search)
	if err != nil {
		return nil, fmt.Errorf("organization service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("organization service: list: count: %w", err)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	data := make([]dto.OrganizationResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	return &dto.OrganizationListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *OrganizationService) Update(ctx context.Context, id string, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("organization service: update: find: %w", err)
	}
	e.Name = req.Name
	e.Domain = req.Domain
	e.IsActive = req.IsActive

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("organization service: update: save: %w", err)
	}

	// Re-fetch to get updated timestamps
	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return s.toResponse(e), nil
	}
	return s.toResponse(updated), nil
}

func (s *OrganizationService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("organization service: delete: %w", err)
	}
	return nil
}

func (s *OrganizationService) toResponse(e *entity.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		Id:        e.Id,
		Name:      e.Name,
		Domain:    e.Domain,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
