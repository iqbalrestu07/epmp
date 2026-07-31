package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/iam/application/dto"
	"github.com/epmp/backend/internal/modules/iam/domain/entity"
	"github.com/epmp/backend/internal/modules/iam/domain/repository"
	"github.com/epmp/backend/internal/shared"
)

// RoleService handles role & permission management use-cases.
type RoleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService creates a new RoleService.
func NewRoleService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *RoleService) Create(ctx context.Context, req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	if _, err := s.roleRepo.FindByName(ctx, req.Name); err == nil {
		return nil, shared.ErrAlreadyExists
	}

	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.roleRepo.Save(ctx, role); err != nil {
		return nil, fmt.Errorf("role service: create: %w", err)
	}

	if len(req.PermissionIDs) > 0 {
		_ = s.roleRepo.SetPermissions(ctx, role.ID, req.PermissionIDs)
	}

	return s.toResponse(ctx, role)
}

func (s *RoleService) GetByID(ctx context.Context, id string) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, shared.ErrNotFound
	}
	return s.toResponse(ctx, role)
}

func (s *RoleService) List(ctx context.Context) ([]*dto.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("role service: list: %w", err)
	}

	resp := make([]*dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		rr, _ := s.toResponse(ctx, r)
		resp = append(resp, rr)
	}
	return resp, nil
}

func (s *RoleService) Update(ctx context.Context, id string, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, shared.ErrNotFound
	}
	if role.IsSystem {
		return nil, shared.NewDomainError("SYSTEM_ROLE", "system roles cannot be modified")
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := s.roleRepo.Save(ctx, role); err != nil {
		return nil, fmt.Errorf("role service: update: %w", err)
	}
	return s.toResponse(ctx, role)
}

func (s *RoleService) SetPermissions(ctx context.Context, id string, req *dto.SetPermissionsRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, shared.ErrNotFound
	}
	if role.IsSystem {
		return nil, shared.NewDomainError("SYSTEM_ROLE", "system roles cannot be modified")
	}

	if err := s.roleRepo.SetPermissions(ctx, id, req.PermissionIDs); err != nil {
		return nil, fmt.Errorf("role service: set permissions: %w", err)
	}
	return s.toResponse(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id string) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return shared.ErrNotFound
	}
	if role.IsSystem {
		return shared.NewDomainError("SYSTEM_ROLE", "system roles cannot be deleted")
	}
	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) ListPermissions(ctx context.Context) ([]*dto.PermissionResponse, error) {
	perms, err := s.permissionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*dto.PermissionResponse, 0, len(perms))
	for _, p := range perms {
		resp = append(resp, &dto.PermissionResponse{
			ID:          p.ID,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
			Key:         p.Key(),
		})
	}
	return resp, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (s *RoleService) toResponse(ctx context.Context, role *entity.Role) (*dto.RoleResponse, error) {
	perms, _ := s.roleRepo.FindPermissionsByRoleID(ctx, role.ID)
	permResps := make([]dto.PermissionResponse, 0, len(perms))
	for _, p := range perms {
		permResps = append(permResps, dto.PermissionResponse{
			ID:          p.ID,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
			Key:         p.Key(),
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		Permissions: permResps,
	}, nil
}
