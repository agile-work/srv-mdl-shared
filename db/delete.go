package db

import (
	"net/http"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Remove object data from the database
func Remove(r *http.Request, scope, table string, conditions builder.Builder) *module.Response {
	response := &module.Response{
		Code: http.StatusOK,
	}

	err := db.DeleteStruct(table, conditions)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorDeletingData, scope, err.Error()))

		return response
	}

	return response
}
