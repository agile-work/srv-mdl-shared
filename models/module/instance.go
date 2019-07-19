package module

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
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

// Add insert a new instance to serve request for this module
func (i *Instance) Add(trs *db.Transaction, moduleCode string) error {
	mdl := Module{Code: moduleCode}
	if err := mdl.Load(); err != nil {
		return customerror.New(http.StatusBadRequest, "load module", err.Error())
	}

	i.ID = db.UUID()

	mdl.Definitions.Instances = append(mdl.Definitions.Instances, *i)

	if err := db.UpdateJSONAttributeTx(trs.Tx, constants.TableCoreModules, "definitions", "{instances}", mdl.Definitions.Instances, builder.Equal("code", moduleCode)); err != nil {
		return err
	}

	return nil
}

// Update changes attributes for the module instance
func (i *Instance) Update(trs *db.Transaction, moduleCode, instanceID string, columns map[string]interface{}) error {
	mdl := Module{Code: moduleCode}
	if err := mdl.Load(); err != nil {
		return customerror.New(http.StatusBadRequest, "load module", err.Error())
	}

	index := getIndexByID(instanceID, mdl.Definitions.Instances)
	if index == -1 {
		return customerror.New(http.StatusNotFound, "update", "instance not found")
	}

	for col, value := range columns {
		path := fmt.Sprintf("{instances, %d, %s}", index, col)
		if err := db.UpdateJSONAttributeTx(trs.Tx, constants.TableCoreModules, "definitions", path, value, builder.Equal("code", moduleCode)); err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a instance from this module
func (i *Instance) Delete(trs *db.Transaction, moduleCode, instanceID string) error {
	mdl := Module{Code: moduleCode}
	if err := mdl.Load(); err != nil {
		return customerror.New(http.StatusBadRequest, "load module", err.Error())
	}

	if len(mdl.Definitions.Instances) <= 1 {
		return customerror.New(http.StatusForbidden, "delete", "at least one instance is required")
	}

	index := getIndexByID(instanceID, mdl.Definitions.Instances)
	if index == -1 {
		return customerror.New(http.StatusNotFound, "delete", "instance not found")
	}

	mdl.Definitions.Instances = append(mdl.Definitions.Instances[:index], mdl.Definitions.Instances[index+1:]...)

	if err := db.UpdateJSONAttributeTx(trs.Tx, constants.TableCoreModules, "definitions", "{instances}", mdl.Definitions.Instances, builder.Equal("code", moduleCode)); err != nil {
		return err
	}

	return nil
}

func getIndexByID(id string, instances []Instance) int {
	for index, instance := range instances {
		if instance.ID == id {
			return index
		}
	}
	return -1
}
