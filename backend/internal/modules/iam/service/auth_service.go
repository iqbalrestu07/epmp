package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/epmp/backend/internal/modules/iam/dto"
	"github.com/epmp/backend/internal/modules/iam/entity"
	"github.com/epmp/backend/internal/modules/iam/repository"
	"github.com/epmp/backend/internal/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication use-cases.
type AuthService struct {
	userRepo         repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
	roleRepo         repository.RoleRepository
	jwtSecret        string
	accessTTL        time.Duration
	refreshTTL       time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	roleRepo repository.RoleRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		userRoleRepo:     userRoleRepo,
		refreshTokenRepo: refreshTokenRepo,
		roleRepo:         roleRepo,
		jwtSecret:        jwtSecret,
		accessTTL:        15 * time.Minute,
		refreshTTL:       7 * 24 * time.Hour,
	}
}

// Register creates the first user (admin) in the system.
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	// Check email uniqueness
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errs.ErrAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth service: register: hash password: %w", err)
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		IsActive:     true,
	}
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("auth service: register: save user: %w", err)
	}

	// Assign super_admin role to the first registered user
	superAdminRole, err := s.roleRepo.FindByName(ctx, "super_admin")
	if err == nil {
		_ = s.userRoleRepo.AssignRoles(ctx, user.ID, []string{superAdminRole.ID})
	}

	return s.buildTokenResponse(ctx, user)
}

// Login authenticates a user and returns JWT tokens.
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	if !user.IsActive {
		return nil, errs.NewDomainError("ACCOUNT_INACTIVE", "account is inactive")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errs.NewDomainError("INVALID_CREDENTIALS", "invalid email or password")
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return s.buildTokenResponse(ctx, user)
}

// Refresh validates a refresh token and issues new tokens.
func (s *AuthService) Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.TokenResponse, error) {
	hash := hashToken(req.RefreshToken)
	rt, err := s.refreshTokenRepo.FindByTokenHash(ctx, hash)
	if err != nil || !rt.IsValid() {
		return nil, errs.NewDomainError("INVALID_TOKEN", "refresh token is invalid or expired")
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil || !user.IsActive {
		return nil, errs.ErrNotFound
	}

	// Rotate: revoke old, issue new
	_ = s.refreshTokenRepo.Revoke(ctx, rt.ID)

	return s.buildTokenResponse(ctx, user)
}

// Logout revokes the given refresh token.
func (s *AuthService) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	hash := hashToken(req.RefreshToken)
	rt, err := s.refreshTokenRepo.FindByTokenHash(ctx, hash)
	if err != nil {
		return nil // idempotent
	}
	return s.refreshTokenRepo.Revoke(ctx, rt.ID)
}

// Me returns the current user profile with permissions.
func (s *AuthService) Me(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errs.ErrNotFound
	}
	return s.buildUserResponse(ctx, user)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (s *AuthService) buildTokenResponse(ctx context.Context, user *entity.User) (*dto.TokenResponse, error) {
	userResp, err := s.buildUserResponse(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := GenerateAccessToken(user.ID, user.Email, s.jwtSecret, s.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("auth service: generate access token: %w", err)
	}

	rawRefresh, err := generateOpaqueToken()
	if err != nil {
		return nil, fmt.Errorf("auth service: generate refresh token: %w", err)
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, fmt.Errorf("auth service: save refresh token: %w", err)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTTL.Seconds()),
		User:         *userResp,
	}, nil
}

func (s *AuthService) buildUserResponse(ctx context.Context, user *entity.User) (*dto.UserResponse, error) {
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

// hashToken returns a SHA-256 hex digest of the token (safe to store in DB).
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// generateOpaqueToken creates a 32-byte random hex token.
func generateOpaqueToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
