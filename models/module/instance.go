package module

import (
	"time"
)

// Instance defines each service instance for this module
type Instance struct {
	ID        string    `json:"id"`
	Host      string    `json:"host" validate:"required"`
	Port      int       `json:"port" validate:"required"`
	CreatedBy string    `json:"created_by" sql:"created_by"`
	CreatedAt time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt time.Time `json:"updated_at" sql:"updated_at"`
}
