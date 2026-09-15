package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/tenantdocument/dto"
	"github.com/epmp/backend/internal/modules/tenantdocument/entity"
	"github.com/epmp/backend/internal/modules/tenantdocument/repository"
)

// TenantDocumentService implements the application layer for TenantDocument.
type TenantDocumentService struct {
	repo repository.TenantDocumentRepository
}

// NewTenantDocumentService creates a new TenantDocumentService.
func NewTenantDocumentService(repo repository.TenantDocumentRepository) *TenantDocumentService {
	return &TenantDocumentService{repo: repo}
}

func (s *TenantDocumentService) Create(ctx context.Context, orgID string, req *dto.CreateTenantDocumentRequest) (*dto.TenantDocumentResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("tenantdocument service: create: organization ID required")
	}
	e := entity.NewTenantDocument()
	e.OrganizationId = orgID
	e.TenantId = req.TenantId
	e.DocumentType = req.DocumentType
	e.FileUrl = req.FileUrl

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantdocument service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantDocumentService) GetByID(ctx context.Context, id, orgID string) (*dto.TenantDocumentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *TenantDocumentService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.TenantDocumentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument service: list: count: %w", err)
	}

	data := make([]dto.TenantDocumentResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.TenantDocumentListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *TenantDocumentService) Update(ctx context.Context, id, orgID string, req *dto.UpdateTenantDocumentRequest) (*dto.TenantDocumentResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument service: update: find: %w", err)
	}
	e.TenantId = req.TenantId
	e.DocumentType = req.DocumentType
	e.FileUrl = req.FileUrl

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("tenantdocument service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *TenantDocumentService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("tenantdocument service: delete: %w", err)
	}
	return nil
}

func (s *TenantDocumentService) toResponse(e *entity.TenantDocument) *dto.TenantDocumentResponse {
	return &dto.TenantDocumentResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		TenantId:       e.TenantId,
		DocumentType:   e.DocumentType,
		FileUrl:        e.FileUrl,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
