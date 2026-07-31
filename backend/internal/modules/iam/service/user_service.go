package service

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/iam/dto"
	"github.com/epmp/backend/internal/modules/iam/entity"
	"github.com/epmp/backend/internal/modules/iam/repository"
	"github.com/epmp/backend/internal/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user management use-cases.
type UserService struct {
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
}

// NewUserService creates a new UserService.
func NewUserService(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errs.ErrAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("user service: create: hash: %w", err)
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		IsActive:     true,
	}
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("user service: create: save: %w", err)
	}

	if len(req.RoleIDs) > 0 {
		_ = s.userRoleRepo.AssignRoles(ctx, user.ID, req.RoleIDs)
	}

	return s.toResponse(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	return s.toResponse(ctx, user)
}

func (s *UserService) List(ctx context.Context, page, perPage int) (*dto.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	users, err := s.userRepo.FindAll(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("user service: list: %w", err)
	}

	total, _ := s.userRepo.CountAll(ctx)

	data := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp, _ := s.toResponse(ctx, u)
		if resp != nil {
			data = append(data, *resp)
		}
	}

	return &dto.UserListResponse{
		Data:    data,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

func (s *UserService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.ErrNotFound
	}

	user.Name = req.Name
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	return s.toResponse(ctx, user)
}

func (s *UserService) AssignRoles(ctx context.Context, id string, req *dto.AssignRolesRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errs.ErrNotFound
	}

	if err := s.userRoleRepo.AssignRoles(ctx, id, req.RoleIDs); err != nil {
		return nil, fmt.Errorf("user service: assign roles: %w", err)
	}

	return s.toResponse(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		return errs.ErrNotFound
	}
	return s.userRepo.Delete(ctx, id)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (s *UserService) toResponse(ctx context.Context, user *entity.User) (*dto.UserResponse, error) {
	roles, _ := s.userRoleRepo.FindRolesByUserID(ctx, user.ID)
	perms, _ := s.userRoleRepo.FindPermissionsByUserID(ctx, user.ID)

	roleResps := make([]dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		roleResps = append(roleResps, dto.RoleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			IsSystem:    r.IsSystem,
		})
	}

	permKeys := make([]string, 0, len(perms))
	for _, p := range perms {
		permKeys = append(permKeys, p.Key())
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		IsActive:    user.IsActive,
		Roles:       roleResps,
		Permissions: permKeys,
	}, nil
}
