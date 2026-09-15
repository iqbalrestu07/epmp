package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/workorder/dto"
	"github.com/epmp/backend/internal/modules/workorder/entity"
	"github.com/epmp/backend/internal/modules/workorder/repository"
)

// WorkOrderService implements the application layer for WorkOrder.
type WorkOrderService struct {
	repo repository.WorkOrderRepository
}

// NewWorkOrderService creates a new WorkOrderService.
func NewWorkOrderService(repo repository.WorkOrderRepository) *WorkOrderService {
	return &WorkOrderService{repo: repo}
}

func (s *WorkOrderService) Create(ctx context.Context, orgID string, req *dto.CreateWorkOrderRequest) (*dto.WorkOrderResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("workorder service: create: organization ID required")
	}
	e := entity.NewWorkOrder()
	e.OrganizationId = orgID
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Description = req.Description
	e.Status = req.Status
	e.Priority = req.Priority

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("workorder service: create: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *WorkOrderService) GetByID(ctx context.Context, id, orgID string) (*dto.WorkOrderResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("workorder service: get by id: %w", err)
	}
	return s.toResponse(e), nil
}

func (s *WorkOrderService) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.WorkOrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("workorder service: list: %w", err)
	}

	total, err := s.repo.Count(ctx, search, orgID)
	if err != nil {
		return nil, fmt.Errorf("workorder service: list: count: %w", err)
	}

	data := make([]dto.WorkOrderResponse, 0, len(items))
	for _, e := range items {
		data = append(data, *s.toResponse(e))
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.WorkOrderListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *WorkOrderService) Update(ctx context.Context, id, orgID string, req *dto.UpdateWorkOrderRequest) (*dto.WorkOrderResponse, error) {
	e, err := s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("workorder service: update: find: %w", err)
	}
	e.PropertyId = req.PropertyId
	e.RoomId = req.RoomId
	e.Description = req.Description
	e.Status = req.Status
	e.Priority = req.Priority

	if err := s.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("workorder service: update: save: %w", err)
	}

	e, err = s.repo.FindByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("workorder service: update: refetch: %w", err)
	}

	return s.toResponse(e), nil
}

func (s *WorkOrderService) Delete(ctx context.Context, id, orgID string) error {
	if err := s.repo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("workorder service: delete: %w", err)
	}
	return nil
}

func (s *WorkOrderService) toResponse(e *entity.WorkOrder) *dto.WorkOrderResponse {
	return &dto.WorkOrderResponse{
		OrganizationId: e.OrganizationId,
		Id:             e.Id,
		PropertyId:     e.PropertyId,
		RoomId:         e.RoomId,
		Description:    e.Description,
		Status:         e.Status,
		Priority:       e.Priority,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
