package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/technician/dto"
	"github.com/epmp/backend/internal/modules/technician/entity"
	"github.com/epmp/backend/internal/modules/technician/repository"
)

// TechnicianService implements the application layer for Technician.
type TechnicianService struct {
	repo repository.TechnicianRepository
}

// NewTechnicianService creates a new TechnicianService.
func NewTechnicianService(repo repository.TechnicianRepository) *TechnicianService {
	return &TechnicianService{repo: repo}
}

func (s *TechnicianService) Create(ctx context.Context, orgID string, req *dto.CreateTechnicianRequest) (*dto.TechnicianResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("technician service: create: organization ID required")
	}
	e := entity.NewTechnician()
	e.OrganizationId = orgID
	e.Name = req.Name
	e.Phone = req.Phone
	e.Specialty = req.Specialty

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("technician service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TechnicianService) GetByID(ctx context.Context, id, orgID string) (*dto.TechnicianResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("technician service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *TechnicianService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.TechnicianListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("technician service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("technician service: list: count: %w", err)
	}

	data := make([]dto.TechnicianResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.TechnicianListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TechnicianService) Update(ctx context.Context, id, orgID string, req *dto.UpdateTechnicianRequest) (*dto.TechnicianResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("technician service: update: find: %w", err)
	}
	e.Name = req.Name
	e.Phone = req.Phone
	e.Specialty = req.Specialty

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("technician service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("technician service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TechnicianService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("technician service: delete: %w", err)
	}
	return nil
}

func (s *TechnicianService) toResponse(e *entity.Technician) *dto.TechnicianResponse {
	return &dto.TechnicianResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		Name: e.Name,
		Phone: e.Phone,
		Specialty: e.Specialty,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
