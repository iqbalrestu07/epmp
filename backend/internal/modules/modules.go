package modules

import (
	"github.com/epmp/backend/internal/modules/building"
	"github.com/epmp/backend/internal/modules/zone"
	"github.com/epmp/backend/internal/modules/bed"
	"github.com/epmp/backend/internal/modules/facility"
	"github.com/epmp/backend/internal/modules/roomtype"
	"github.com/epmp/backend/internal/modules/tenantidentity"
	"github.com/epmp/backend/internal/modules/tenantcontact"
	"github.com/epmp/backend/internal/modules/tenantdocument"
	"github.com/epmp/backend/internal/modules/asset"
	"github.com/epmp/backend/internal/modules/assetassignment"
	"github.com/epmp/backend/internal/modules/assetinspection"
	"github.com/epmp/backend/internal/modules/workorder"
	"github.com/epmp/backend/internal/modules/technician"
	"github.com/epmp/backend/internal/modules/supplier"

	"github.com/epmp/backend/internal/modules/adjustment"
	"github.com/epmp/backend/internal/modules/billing"
	"github.com/epmp/backend/internal/modules/charge"
	"github.com/epmp/backend/internal/modules/contract"
	"github.com/epmp/backend/internal/modules/deposit"
	"github.com/epmp/backend/internal/modules/iam"
	"github.com/epmp/backend/internal/modules/occupancy"
	"github.com/epmp/backend/internal/modules/organization"
	"github.com/epmp/backend/internal/modules/payment"
	"github.com/epmp/backend/internal/modules/penalty"
	"github.com/epmp/backend/internal/modules/property"
	"github.com/epmp/backend/internal/modules/refund"
	"github.com/epmp/backend/internal/modules/reservation"
	"github.com/epmp/backend/internal/modules/room"
	"github.com/epmp/backend/internal/modules/tenant"
	mw "github.com/epmp/backend/internal/pkg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Register registers all modules to the Echo router.
func Register(e *echo.Echo, db *pgxpool.Pool, log zerolog.Logger, jwtSecret string) error {
	v1 := e.Group("/api/v1")

	// IAM module registers its own public and protected routes.
	iamMod := iam.NewModule(db, log, jwtSecret)
	iamMod.RegisterRoutes(v1)

	// Resource modules — all protected by JWT auth middleware.
	protected := v1.Group("", mw.AuthRequired(jwtSecret))
	property.NewModule(db, log).RegisterRoutes(protected)
	tenant.NewModule(db, log).RegisterRoutes(protected)
	room.NewModule(db, log).RegisterRoutes(protected)
	reservation.NewModule(db, log).RegisterRoutes(protected)
	contract.NewModule(db, log).RegisterRoutes(protected)
	organization.NewModule(db, log).RegisterRoutes(protected)
	occupancy.NewModule(db, log).RegisterRoutes(protected)
	billing.NewModule(db, log).RegisterRoutes(protected)
	payment.NewModule(db, log).RegisterRoutes(protected)
	deposit.NewModule(db, log).RegisterRoutes(protected)
	charge.NewModule(db, log).RegisterRoutes(protected)
	refund.NewModule(db, log).RegisterRoutes(protected)
	adjustment.NewModule(db, log).RegisterRoutes(protected)
	penalty.NewModule(db, log).RegisterRoutes(protected)

		building.NewModule(db, log).RegisterRoutes(protected)
	zone.NewModule(db, log).RegisterRoutes(protected)
	bed.NewModule(db, log).RegisterRoutes(protected)
	facility.NewModule(db, log).RegisterRoutes(protected)
	roomtype.NewModule(db, log).RegisterRoutes(protected)
	tenantidentity.NewModule(db, log).RegisterRoutes(protected)
	tenantcontact.NewModule(db, log).RegisterRoutes(protected)
	tenantdocument.NewModule(db, log).RegisterRoutes(protected)
	asset.NewModule(db, log).RegisterRoutes(protected)
	assetassignment.NewModule(db, log).RegisterRoutes(protected)
	assetinspection.NewModule(db, log).RegisterRoutes(protected)
	workorder.NewModule(db, log).RegisterRoutes(protected)
	technician.NewModule(db, log).RegisterRoutes(protected)
	supplier.NewModule(db, log).RegisterRoutes(protected)

	return nil
}
