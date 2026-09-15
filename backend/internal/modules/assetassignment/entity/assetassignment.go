package entity

import "time"

// AssetAssignment is the domain entity for assetassignment.
type AssetAssignment struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	AssetId        string     `json:"asset_id"`
	RoomId         string     `json:"room_id"`
	AssignedDate   time.Time  `json:"assigned_date"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewAssetAssignment creates a new AssetAssignment instance.
func NewAssetAssignment() *AssetAssignment {
	return &AssetAssignment{}
}
