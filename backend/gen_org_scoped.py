#!/usr/bin/env python3
"""Generate org-scoped filtering for all remaining modules."""
import os

BASE = "/Users/macbookpro/pjc/personal/epmp/backend/internal/modules"

# Module configs: (module_dir, entity_name, table_name, fields)
# fields: list of (go_field, json_name, db_column, is_time)
MODULES = [
    {
        "dir": "reservation",
        "entity": "Reservation",
        "table": "reservations",
        "search_fields": ["status"],
        "fields": [
            ("TenantId", "tenant_id", "tenant_id", False),
            ("PropertyId", "property_id", "property_id", False),
            ("RoomId", "room_id", "room_id", False),
            ("Status", "status", "status", False),
            ("CheckInDate", "check_in_date", "check_in_date", True),
            ("CheckOutDate", "check_out_date", "check_out_date", True),
            ("BookingFee", "booking_fee", "booking_fee", False),
            ("Notes", "notes", "notes", False),
        ],
    },
    {
        "dir": "contract",
        "entity": "Contract",
        "table": "contracts",
        "search_fields": ["status"],
        "fields": [
            ("ReservationId", "reservation_id", "reservation_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("PropertyId", "property_id", "property_id", False),
            ("RoomId", "room_id", "room_id", False),
            ("Status", "status", "status", False),
            ("StartDate", "start_date", "start_date", True),
            ("EndDate", "end_date", "end_date", True),
            ("MonthlyRent", "monthly_rent", "monthly_rent", False),
            ("DepositAmount", "deposit_amount", "deposit_amount", False),
            ("Terms", "terms", "terms", False),
        ],
    },
    {
        "dir": "occupancy",
        "entity": "Occupancy",
        "table": "occupancies",
        "search_fields": ["status"],
        "fields": [
            ("ContractId", "contract_id", "contract_id", False),
            ("RoomId", "room_id", "room_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("Status", "status", "status", False),
            ("CheckInTime", "check_in_time", "check_in_time", True),
            ("CheckOutTime", "check_out_time", "check_out_time", True),
            ("Notes", "notes", "notes", False),
        ],
    },
    {
        "dir": "billing",
        "entity": "Invoice",
        "table": "invoices",
        "search_fields": ["status"],
        "fields": [
            ("ContractId", "contract_id", "contract_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("Amount", "amount", "amount", False),
            ("Status", "status", "status", False),
            ("DueDate", "due_date", "due_date", True),
            ("PaidDate", "paid_date", "paid_date", True),
            ("PaymentMethod", "payment_method", "payment_method", False),
            ("Notes", "notes", "notes", False),
        ],
    },
    {
        "dir": "payment",
        "entity": "Payment",
        "table": "payments",
        "search_fields": ["status", "payment_method"],
        "fields": [
            ("InvoiceId", "invoice_id", "invoice_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("Amount", "amount", "amount", False),
            ("PaymentDate", "payment_date", "payment_date", True),
            ("PaymentMethod", "payment_method", "payment_method", False),
            ("Status", "status", "status", False),
            ("ReferenceNumber", "reference_number", "reference_number", False),
        ],
    },
    {
        "dir": "deposit",
        "entity": "Deposit",
        "table": "deposits",
        "search_fields": ["status"],
        "fields": [
            ("ContractId", "contract_id", "contract_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("Amount", "amount", "amount", False),
            ("Status", "status", "status", False),
            ("CollectionDate", "collection_date", "collection_date", True),
            ("RefundDate", "refund_date", "refund_date", True),
            ("Notes", "notes", "notes", False),
        ],
    },
    {
        "dir": "charge",
        "entity": "Charge",
        "table": "charges",
        "search_fields": ["charge_type", "status"],
        "fields": [
            ("ContractId", "contract_id", "contract_id", False),
            ("InvoiceId", "invoice_id", "invoice_id", False),
            ("ChargeType", "charge_type", "charge_type", False),
            ("Amount", "amount", "amount", False),
            ("Status", "status", "status", False),
            ("ChargeDate", "charge_date", "charge_date", True),
            ("Notes", "notes", "notes", False),
        ],
    },
    {
        "dir": "refund",
        "entity": "Refund",
        "table": "refunds",
        "search_fields": ["status"],
        "fields": [
            ("PaymentId", "payment_id", "payment_id", False),
            ("TenantId", "tenant_id", "tenant_id", False),
            ("Amount", "amount", "amount", False),
            ("Status", "status", "status", False),
            ("RefundDate", "refund_date", "refund_date", True),
            ("Reason", "reason", "reason", False),
        ],
    },
    {
        "dir": "adjustment",
        "entity": "Adjustment",
        "table": "adjustments",
        "search_fields": ["adjustment_type"],
        "fields": [
            ("InvoiceId", "invoice_id", "invoice_id", False),
            ("AdjustmentType", "adjustment_type", "adjustment_type", False),
            ("Amount", "amount", "amount", False),
            ("AdjustmentDate", "adjustment_date", "adjustment_date", True),
            ("Reason", "reason", "reason", False),
        ],
    },
    {
        "dir": "penalty",
        "entity": "Penalty",
        "table": "penalties",
        "search_fields": ["status"],
        "fields": [
            ("InvoiceId", "invoice_id", "invoice_id", False),
            ("Amount", "amount", "amount", False),
            ("Status", "status", "status", False),
            ("PenaltyDate", "penalty_date", "penalty_date", True),
            ("Description", "description", "description", False),
        ],
    },
]


def gen_entity(mod):
    e = mod["entity"]
    fields_str = "\n".join(
        f'\t{f[0]:<20} {("time.Time" if f[3] else "string"):<12} `json:"{f[1]}"`'
        for f in mod["fields"]
    )
    return f"""package entity

import "time"

// {e} is the domain entity for {mod["dir"]}.
type {e} struct {{
\tOrganizationId string     `json:"organization_id"`
\tId             string     `json:"id"`
{fields_str}
\tDeletedAt      *time.Time `json:"deleted_at,omitempty"`
\tCreatedAt      time.Time  `json:"created_at"`
\tUpdatedAt      time.Time  `json:"updated_at"`
}}

// New{e} creates a new {e} instance.
func New{e}() *{e} {{
\treturn &{e}{{}}
}}
"""


def gen_dto(mod):
    e = mod["entity"]
    fields = mod["fields"]
    create_fields = "\n".join(
        f'\t{f[0]:<20} {("time.Time" if f[3] else "string"):<12} `json:"{f[1]}"`'
        if f[3] or f[0] != "Amount" and not f[0].endswith("Fee") and not f[0].endswith("Rent") and not f[0].endswith("Amount")
        else f'\t{f[0]:<20} float64    `json:"{f[1]}"`'
        for f in fields
    )
    # Fix: Amount, BookingFee, MonthlyRent, DepositAmount are float64
    create_fields_lines = []
    for f in fields:
        go_type = "time.Time" if f[3] else "string"
        if f[0] in ("Amount", "BookingFee", "MonthlyRent", "DepositAmount"):
            go_type = "float64"
        create_fields_lines.append(f'\t{f[0]:<20} {go_type:<12} `json:"{f[1]}"`')
    create_fields = "\n".join(create_fields_lines)

    resp_fields_lines = []
    for f in fields:
        go_type = "time.Time" if f[3] else "string"
        if f[0] in ("Amount", "BookingFee", "MonthlyRent", "DepositAmount"):
            go_type = "float64"
        resp_fields_lines.append(f'\t{f[0]:<20} {go_type:<12} `json:"{f[1]}"`')
    resp_fields = "\n".join(resp_fields_lines)

    return f"""package dto

import "time"

// Create{e}Request is the DTO for creating a {e}.
type Create{e}Request struct {{
{create_fields}
}}

// Update{e}Request is the DTO for updating a {e}.
type Update{e}Request struct {{
{create_fields}
}}

// {e}Response is the DTO for returning a {e}.
type {e}Response struct {{
\tOrganizationId string    `json:"organization_id"`
\tId             string    `json:"id"`
{resp_fields}
\tCreatedAt      time.Time `json:"created_at"`
\tUpdatedAt      time.Time `json:"updated_at"`
}}

// {e}ListResponse is the DTO for a paginated list of {e}.
type {e}ListResponse struct {{
\tData       []{e}Response `json:"data"`
\tTotal      int64         `json:"total"`
\tPage       int           `json:"page"`
\tPerPage    int           `json:"per_page"`
\tTotalPages int           `json:"total_pages"`
}}
"""


def gen_repo_interface(mod):
    e = mod["entity"]
    return f"""package repository

import (
\t"context"

\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/entity"
)

// {e}Repository is the contract for persisting {e} aggregates.
type {e}Repository interface {{
\t// Save persists a {e} entity.
\tSave(ctx context.Context, e *entity.{e}) error

\t// FindByID retrieves a {e} by its primary key within an organization.
\tFindByID(ctx context.Context, id, orgID string) (*entity.{e}, error)

\t// FindAll retrieves a paginated list of {e} entities within an organization.
\tFindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.{e}, error)

\t// Count returns the total number of {e} entities matching the filter.
\tCount(ctx context.Context, search, orgID string) (int64, error)

\t// Delete removes a {e} by its primary key.
\tDelete(ctx context.Context, id, orgID string) error
}}
"""


def gen_repo_impl(mod):
    e = mod["entity"]
    table = mod["table"]
    fields = mod["fields"]
    search_fields = mod["search_fields"]

    # Column names for INSERT/SELECT
    db_cols = [f[2] for f in fields]
    all_cols = ["organization_id", "id"] + db_cols + ["deleted_at", "created_at", "updated_at"]
    select_cols = ", ".join(all_cols)

    # INSERT columns (without id, deleted_at, created_at, updated_at)
    insert_cols = ", ".join(["organization_id"] + db_cols)
    insert_placeholders = ", ".join(f"${{i}}" for i in range(1, len(db_cols) + 2))
    insert_args = ", ".join(["e.OrganizationId"] + [f"e.{f[0]}" for f in fields])

    # UPDATE columns (without organization_id)
    update_cols = db_cols
    update_set = ", ".join(f"{update_cols[i]}=${{i+1}}" for i in range(len(update_cols)))
    update_args = ", ".join([f"e.{f[0]}" for f in fields] + ["e.Id", "e.OrganizationId"])

    # Scan fields for FindByID
    scan_fields = ", ".join(["&e.OrganizationId", "&e.Id"] + [f"&e.{f[0]}" for f in fields] + ["&e.DeletedAt", "&e.CreatedAt", "&e.UpdatedAt"])

    # Search condition
    if len(search_fields) == 1:
        search_cond = f'{search_fields[0]} ILIKE ${{argIdx}}'
    else:
        search_cond = " OR ".join(f'{sf} ILIKE ${{argIdx}}' for sf in search_fields)
        search_cond = f"({search_cond})"

    return f"""package repository

import (
\t"context"
\t"fmt"

\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/entity"

\t"github.com/jackc/pgx/v5/pgxpool"
)

// {e}RepositoryImpl implements {e}Repository using PostgreSQL.
type {e}RepositoryImpl struct {{
\tdb *pgxpool.Pool
}}

// New{e}RepositoryImpl creates a new {e}RepositoryImpl.
func New{e}RepositoryImpl(db *pgxpool.Pool) *{e}RepositoryImpl {{
\treturn &{e}RepositoryImpl{{db: db}}
}}

// Ensure {e}RepositoryImpl implements domain repository interface.
var _ {e}Repository = (*{e}RepositoryImpl)(nil)

func (r *{e}RepositoryImpl) Save(ctx context.Context, e *entity.{e}) error {{
\tif e.Id == "" {{
\t\terr := r.db.QueryRow(ctx, `
\t\t\tINSERT INTO {table} (organization_id, {", ".join(db_cols)})
\t\t\tVALUES ($1, {", ".join(f"${{i+2}}" for i in range(len(db_cols)))})
\t\t\tRETURNING id, created_at, updated_at`,
\t\t\te.OrganizationId, {", ".join(f"e.{f[0]}" for f in fields)},
\t\t).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
\t\treturn err
\t}}
\terr := r.db.QueryRow(ctx, `
\t\tUPDATE {table}
\t\tSET    {", ".join(f"{db_cols[i]}=${{i+1}}" for i in range(len(db_cols)))}
\t\tWHERE  id=${len(db_cols)+1} AND organization_id=${len(db_cols)+2} AND deleted_at IS NULL
\t\tRETURNING updated_at`,
\t\t{", ".join(f"e.{f[0]}" for f in fields)}, e.Id, e.OrganizationId,
\t).Scan(&e.UpdatedAt)
\treturn err
}}

func (r *{e}RepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.{e}, error) {{
\te := &entity.{e}{{}}
\terr := r.db.QueryRow(ctx, `
\t\tSELECT {select_cols}
\t\tFROM   {table}
\t\tWHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
\t\tid, orgID,
\t).Scan({scan_fields})

\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} repository: find by id: %w", err)
\t}}
\treturn e, nil
}}

func (r *{e}RepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.{e}, error) {{
\tquery := `
\t\tSELECT {select_cols}
\t\tFROM   {table}
\t\tWHERE  deleted_at IS NULL AND organization_id = $1`
\targs := []interface{}{{orgID}}
\targIdx := 2

\tif search != "" {{
\t\tquery += fmt.Sprintf(" AND {search_cond}", argIdx)
\t\targs = append(args, "%"+search+"%")
\t\targIdx++
\t}}

\tquery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
\targs = append(args, limit, offset)

\trows, err := r.db.Query(ctx, query, args...)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} repository: find all: %w", err)
\t}}
\tdefer rows.Close()

\tvar list []*entity.{e}
\tfor rows.Next() {{
\t\te := &entity.{e}{{}}
\t\tif err := rows.Scan({scan_fields}); err != nil {{
\t\t\treturn nil, err
\t\t}}
\t\tlist = append(list, e)
\t}}
\treturn list, rows.Err()
}}

func (r *{e}RepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {{
\tquery := `SELECT COUNT(*) FROM {table} WHERE deleted_at IS NULL AND organization_id = $1`
\targs := []interface{}{{orgID}}

\tif search != "" {{
\t\tquery += fmt.Sprintf(" AND {search_cond}", 2)
\t\targs = append(args, "%"+search+"%")
\t}}

\tvar count int64
\terr := r.db.QueryRow(ctx, query, args...).Scan(&count)
\tif err != nil {{
\t\treturn 0, fmt.Errorf("{mod["dir"]} repository: count: %w", err)
\t}}
\treturn count, nil
}}

func (r *{e}RepositoryImpl) Delete(ctx context.Context, id, orgID string) error {{
\t_, err := r.db.Exec(ctx, `
\t\tUPDATE {table} SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
\treturn err
}}
"""


def gen_service(mod):
    e = mod["entity"]
    fields = mod["fields"]

    create_lines = "\n".join(f"\te.{f[0]} = req.{f[0]}" for f in fields)
    update_lines = "\n".join(f"\te.{f[0]} = req.{f[0]}" for f in fields)
    resp_lines = ", ".join([f"OrganizationId: e.OrganizationId", f"Id: e.Id"] + [f"{f[0]}: e.{f[0]}" for f in fields] + ["CreatedAt: e.CreatedAt", "UpdatedAt: e.UpdatedAt"])

    return f"""package service

import (
\t"context"
\t"fmt"
\t"math"

\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/dto"
\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/entity"
\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/repository"
)

// {e}Service implements the application layer for {e}.
type {e}Service struct {{
\trepo repository.{e}Repository
}}

// New{e}Service creates a new {e}Service.
func New{e}Service(repo repository.{e}Repository) *{e}Service {{
\treturn &{e}Service{{repo: repo}}
}}

func (s *{e}Service) Create(ctx context.Context, orgID string, req *dto.Create{e}Request) (*dto.{e}Response, error) {{
\tif orgID == "" {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: create: organization ID required")
\t}}
\te := entity.New{e}()
\te.OrganizationId = orgID
{create_lines}

\tif err := s.repo.Save(ctx, e); err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: create: %w", err)
\t}}

\treturn s.toResponse(e), nil
}}

func (s *{e}Service) GetByID(ctx context.Context, id, orgID string) (*dto.{e}Response, error) {{
\te, err := s.repo.FindByID(ctx, id, orgID)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: get by id: %w", err)
\t}}
\treturn s.toResponse(e), nil
}}

func (s *{e}Service) List(ctx context.Context, page, perPage int, search, orgID string) (*dto.{e}ListResponse, error) {{
\tif page < 1 {{
\t\tpage = 1
\t}}
\tif perPage < 1 {{
\t\tperPage = 20
\t}}
\toffset := (page - 1) * perPage

\titems, err := s.repo.FindAll(ctx, perPage, offset, search, orgID)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: list: %w", err)
\t}}

\ttotal, err := s.repo.Count(ctx, search, orgID)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: list: count: %w", err)
\t}}

\tdata := make([]dto.{e}Response, 0, len(items))
\tfor _, e := range items {{
\t\tdata = append(data, *s.toResponse(e))
\t}}

\ttotalPages := int(math.Ceil(float64(total) / float64(perPage)))

\treturn &dto.{e}ListResponse{{
\t\tData:       data,
\t\tTotal:      total,
\t\tPage:       page,
\t\tPerPage:    perPage,
\t\tTotalPages: totalPages,
\t}}, nil
}}

func (s *{e}Service) Update(ctx context.Context, id, orgID string, req *dto.Update{e}Request) (*dto.{e}Response, error) {{
\te, err := s.repo.FindByID(ctx, id, orgID)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: update: find: %w", err)
\t}}
{update_lines}

\tif err := s.repo.Save(ctx, e); err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: update: save: %w", err)
\t}}

\te, err = s.repo.FindByID(ctx, id, orgID)
\tif err != nil {{
\t\treturn nil, fmt.Errorf("{mod["dir"]} service: update: refetch: %w", err)
\t}}

\treturn s.toResponse(e), nil
}}

func (s *{e}Service) Delete(ctx context.Context, id, orgID string) error {{
\tif err := s.repo.Delete(ctx, id, orgID); err != nil {{
\t\treturn fmt.Errorf("{mod["dir"]} service: delete: %w", err)
\t}}
\treturn nil
}}

func (s *{e}Service) toResponse(e *entity.{e}) *dto.{e}Response {{
\treturn &dto.{e}Response{{
\t\t{resp_lines},
\t}}
}}
"""


def gen_handler(mod):
    e = mod["entity"]
    return f"""package http

import (
\t"strconv"

\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/dto"
\t"github.com/epmp/backend/internal/modules/{mod["dir"]}/service"
\tmw "github.com/epmp/backend/internal/pkg/middleware"
\t"github.com/epmp/backend/internal/pkg/response"

\t"github.com/labstack/echo/v4"
)

// {e}Handler handles HTTP requests for {e} resources.
type {e}Handler struct {{
\tsvc *service.{e}Service
}}

// New{e}Handler creates a new {e}Handler.
func New{e}Handler(svc *service.{e}Service) *{e}Handler {{
\treturn &{e}Handler{{svc: svc}}
}}

func (h *{e}Handler) Create(c echo.Context) error {{
\torgID := mw.GetOrgID(c)
\tif orgID == "" {{
\t\treturn response.BadRequest(c, "X-Organization-ID header is required")
\t}}

\tvar req dto.Create{e}Request
\tif err := c.Bind(&req); err != nil {{
\t\treturn response.BadRequest(c, "invalid request body")
\t}}

\tresult, err := h.svc.Create(c.Request().Context(), orgID, &req)
\tif err != nil {{
\t\treturn response.InternalError(c, err.Error())
\t}}

\treturn response.Created(c, result)
}}

func (h *{e}Handler) GetByID(c echo.Context) error {{
\tid := c.Param("id")
\torgID := mw.GetOrgID(c)

\tresult, err := h.svc.GetByID(c.Request().Context(), id, orgID)
\tif err != nil {{
\t\treturn response.NotFound(c, "{e} not found")
\t}}

\treturn response.OK(c, result)
}}

func (h *{e}Handler) List(c echo.Context) error {{
\tpage, _ := strconv.Atoi(c.QueryParam("page"))
\tif page == 0 {{
\t\tpage = 1
\t}}
\tperPage, _ := strconv.Atoi(c.QueryParam("per_page"))
\tif perPage == 0 {{
\t\tperPage = 20
\t}}
\tsearch := c.QueryParam("search")
\torgID := mw.GetOrgID(c)

\tresult, err := h.svc.List(c.Request().Context(), page, perPage, search, orgID)
\tif err != nil {{
\t\treturn response.InternalError(c, err.Error())
\t}}

\treturn response.OK(c, result)
}}

func (h *{e}Handler) Update(c echo.Context) error {{
\tid := c.Param("id")
\torgID := mw.GetOrgID(c)

\tvar req dto.Update{e}Request
\tif err := c.Bind(&req); err != nil {{
\t\treturn response.BadRequest(c, "invalid request body")
\t}}

\tresult, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
\tif err != nil {{
\t\treturn response.InternalError(c, err.Error())
\t}}

\treturn response.OK(c, result)
}}

func (h *{e}Handler) Delete(c echo.Context) error {{
\tid := c.Param("id")
\torgID := mw.GetOrgID(c)

\tif err := h.svc.Delete(c.Request().Context(), id, orgID); err != nil {{
\t\treturn response.InternalError(c, err.Error())
\t}}

\treturn response.NoContent(c)
}}
"""


def main():
    for mod in MODULES:
        d = os.path.join(BASE, mod["dir"])
        print(f"Generating {mod['dir']}...")

        # Entity
        with open(os.path.join(d, "entity", f"{mod['dir']}.go"), "w") as f:
            f.write(gen_entity(mod))

        # DTO
        with open(os.path.join(d, "dto", f"{mod['dir']}_dto.go"), "w") as f:
            f.write(gen_dto(mod))

        # Repository interface
        with open(os.path.join(d, "repository", f"{mod['dir']}_repository.go"), "w") as f:
            f.write(gen_repo_interface(mod))

        # Repository impl
        with open(os.path.join(d, "repository", f"{mod['dir']}_repository_impl.go"), "w") as f:
            f.write(gen_repo_impl(mod))

        # Service
        with open(os.path.join(d, "service", f"{mod['dir']}_service.go"), "w") as f:
            f.write(gen_service(mod))

        # Handler
        with open(os.path.join(d, "delivery", "http", f"{mod['dir']}_handler.go"), "w") as f:
            f.write(gen_handler(mod))

    print("Done!")


if __name__ == "__main__":
    main()
