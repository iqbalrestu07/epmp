package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/organization/dto"
	"github.com/epmp/backend/internal/modules/organization/entity"
	"github.com/epmp/backend/internal/modules/organization/repository"
	"github.com/epmp/backend/internal/pkg/uid"
)

// OrganizationService implements the application layer for Organization.
type OrganizationService struct {
	repo repository.OrganizationRepository
}

// NewOrganizationService creates a new OrganizationService.
func NewOrganizationService(repo repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// Create creates a new organization and automatically assigns the creator as owner.
// createdBy must be the authenticated user's ID from JWT context.
func (s *OrganizationService) Create(ctx context.Context, req *dto.CreateOrganizationRequest, createdBy string) (*dto.OrganizationResponse, error) {
	e := entity.NewOrganization()
	e.Id = uid.New()
	e.Name = req.Name
	e.Domain = req.Domain
	e.IsActive = req.IsActive
	e.CreatedBy = createdBy

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("organization service: create: %w", err)
	}

	// Auto-assign creator as owner in organization_members
	member := entity.NewOrganizationMember(e.Id, createdBy, entity.OrgRoleOwner)
	member.Id = uid.New()
	if err := s.repo.SaveMember(ctx, member); err != nil {
		// Non-fatal: log but don't fail the org creation
		// In production, wrap in a transaction
		return nil, fmt.Errorf("organization service: create: assign owner: %w", err)
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

// ListByUser returns all organizations where the given user is a member.
func (s *OrganizationService) ListByUser(ctx context.Context, userID string) ([]*dto.OrganizationResponse, error) {
	items, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("organization service: list by user: %w", err)
	}
	result := make([]*dto.OrganizationResponse, 0, len(items))
	for _, e := range items {
		result = append(result, s.toResponse(e))
	}
	return result, nil
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

func (s *OrganizationService) ListMembers(ctx context.Context, orgID string) ([]*dto.OrganizationMemberResponse, error) {
	return s.repo.FindMembersWithUserByOrgID(ctx, orgID)
}

func (s *OrganizationService) AddMember(ctx context.Context, orgID string, req *dto.AddOrganizationMemberRequest, invitedBy string) (*dto.OrganizationMemberResponse, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	return s.repo.AddMemberByEmail(ctx, orgID, req.Email, req.Name, req.Password, req.Role, invitedBy)
}

func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, userID string) error {
	return s.repo.DeleteMember(ctx, orgID, userID)
}


func (s *OrganizationService) toResponse(e *entity.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		Id:        e.Id,
		Name:      e.Name,
		Domain:    e.Domain,
		IsActive:  e.IsActive,
		CreatedBy: e.CreatedBy,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
