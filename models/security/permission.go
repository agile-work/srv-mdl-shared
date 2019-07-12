package security

import (
	"encoding/json"
	"time"
)

// Permission defines the struct of this object
type Permission struct { // TODO: Rever conceitualmente a permissão
	Type                int              `json:"permission_type"` // 100 (read), 200 (edit)
	StructureType       string           `json:"structure_type"`
	StructureDefinition json.RawMessage  `json:"structure_definition"`
	Users               []PermissionUser `json:"users,omitempty"`
}

// PermissionUser defines the struct to the jsonb field in group
type PermissionUser struct {
	Username  string    `json:"username"`
	CreatedBy string    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
