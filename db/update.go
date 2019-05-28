package db

import (
	"fmt"
	"net/http"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Update object data in the database
func Update(r *http.Request, object interface{}, scope, table string, condition builder.Builder) *module.Response {
	response := GetResponse(r, object, scope)
	if response.Code != http.StatusOK {
		return response
	}

	columns := getColumnsFromBody(r, object)

	err := db.UpdateStruct(table, object, condition, columns...)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update", scope), err.Error()))

		return response
	}

	translationColumns := GetTranslationLanguageCodeColumns(object, columns...)

	if len(translationColumns) > 0 {
		err := UpdateTranslationsFromStruct(table, r.Header.Get("Content-Language"), object, columns...)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update translation", scope), err.Error()))

			return response
		}
	}

	response.Data = object

	return response
}
