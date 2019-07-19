package module

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/feature"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
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

// Register defines the configuration for a new module
func (m *Module) Register(trs *db.Transaction) error {
	translation.SetStructTranslationsLanguage(m, "all")
	m.Status = constants.ModuleStatusRegistered
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreModules, m)
	if err != nil {
		return customerror.New(http.StatusInternalServerError, "module register", err.Error())
	}
	m.ID = id
	return nil
}

// Update updates object data in the database
func (m *Module) Update(trs *db.Transaction, columns []string, translations map[string]string) error {
	opt := &db.Options{Conditions: builder.Equal("code", m.Code)}

	if len(columns) > 0 {
		if err := db.UpdateStructTx(trs.Tx, constants.TableCoreModules, m, opt, strings.Join(columns, ",")); err != nil {
			return customerror.New(http.StatusInternalServerError, "module update", err.Error())
		}
	}

	if len(translations) > 0 {
		statement := builder.Update(constants.TableCoreModules)
		for col, val := range translations {
			statement.JSON(col, translation.FieldsRequestLanguageCode)
			jsonVal, _ := json.Marshal(val)
			statement.Values(jsonVal)
		}
		statement.Where(opt.Conditions)
		if _, err := trs.Query(statement); err != nil {
			return customerror.New(http.StatusInternalServerError, "module update", err.Error())
		}
	}

	return nil
}

// Load returns only one object from the database
func (m *Module) Load() error {
	if err := db.SelectStruct(constants.TableCoreModules, m, &db.Options{
		Conditions: builder.Equal("code", m.Code),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "module load", err.Error())
	}
	return nil
}

// Delete the object from the database
func (m *Module) Delete(trs *db.Transaction) error {
	if err := db.DeleteStructTx(trs.Tx, constants.TableCoreModules, &db.Options{
		Conditions: builder.Equal("code", m.Code),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "module create", err.Error())
	}
	return nil
}

// Modules slice of module
type Modules []Module

// LoadAll defines all instances from the object
func (m *Modules) LoadAll(opt *db.Options) error {
	if err := db.SelectStruct(constants.TableCoreModules, m, opt); err != nil {
		return customerror.New(http.StatusInternalServerError, "modules load", err.Error())
	}
	return nil
}
