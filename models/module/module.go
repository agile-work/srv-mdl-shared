package module

import (
	"net/http"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/feature"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Module define a new module in the application
type Module struct {
	ID          string                  `json:"id" sql:"id" pk:"true"`
	Code        string                  `json:"code" sql:"code" updatable:"false" validate:"required"`
	Prefix      string                  `json:"prefix" sql:"prefix" validate:"required"`
	Status      string                  `json:"status" sql:"status"`
	Name        translation.Translation `json:"name" sql:"name" field:"jsonb" validate:"required"`
	Description translation.Translation `json:"description" sql:"description" field:"jsonb"`
	Version     string                  `json:"version" sql:"version" updatable:"false" validate:"required"`
	Definitions Definition              `json:"definitions,omitempty" sql:"definitions" updatable:"false" field:"jsonb"`
	IsSystem    bool                    `json:"is_system" sql:"is_system"`
	Active      bool                    `json:"active" sql:"active"`
	CreatedBy   string                  `json:"created_by" sql:"created_by"`
	CreatedAt   time.Time               `json:"created_at" sql:"created_at"`
	UpdatedBy   string                  `json:"updated_by" sql:"updated_by"`
	UpdatedAt   time.Time               `json:"updated_at" sql:"updated_at"`
}

// Definition configurations for this module
type Definition struct {
	Instances []Instance                 `json:"instances,omitempty"`
	Features  map[string]feature.Feature `json:"features,omitempty"`
}

// Create defines the configuration for a new module
func (m *Module) Create(trs *db.Transaction) error {
	translation.SetStructTranslationsLanguage(m, "all")
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreModules, m)
	if err != nil {
		return customerror.New(http.StatusInternalServerError, "module register", err.Error())
	}
	m.ID = id
	return nil
}
