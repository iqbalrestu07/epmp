package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/reservation/dto"
	"github.com/epmp/backend/internal/modules/reservation/entity"
	"github.com/epmp/backend/internal/modules/reservation/repository"
)

// ReservationService implements the application layer for Reservation.
type ReservationService struct {
	repo repository.ReservationRepository
}

// NewReservationService creates a new ReservationService.
func NewReservationService(repo repository.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Create(ctx context.Context, orgID string, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("reservation service: create: organization ID required")
	}
	e := entity.NewReservation()
	e.OrganizationId = orgID
	e.TenantId = req.TenantId
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Status = req.Status
	e.CheckInDate = req.CheckInDate
	e.CheckOutDate = req.CheckOutDate
	e.BookingFee = req.BookingFee
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("reservation service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ReservationService) GetByID(ctx context.Context, id, orgID string) (*dto.ReservationResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("reservation service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *ReservationService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.ReservationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("reservation service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("reservation service: list: count: %w", err)
	}

	data := make([]dto.ReservationResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.ReservationListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *ReservationService) Update(ctx context.Context, id, orgID string, req *dto.UpdateReservationRequest) (*dto.ReservationResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("reservation service: update: find: %w", err)
	}
	e.TenantId = req.TenantId
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Status = req.Status
	e.CheckInDate = req.CheckInDate
	e.CheckOutDate = req.CheckOutDate
	e.BookingFee = req.BookingFee
	e.Notes = req.Notes

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("reservation service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("reservation service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *ReservationService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("reservation service: delete: %w", err)
	}
	return nil
}

func (s *ReservationService) toResponse(e *entity.Reservation) *dto.ReservationResponse {
	return &dto.ReservationResponse{
		OrganizationId: e.OrganizationId,
		Id: e.Id,
		TenantId: e.TenantId,
		PropertyId: e.PropertyId,
		RoomId: e.RoomId,
		Status: e.Status,
		CheckInDate: e.CheckInDate,
		CheckOutDate: e.CheckOutDate,
		BookingFee: e.BookingFee,
		Notes: e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
