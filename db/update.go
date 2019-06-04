package db

import (
	"encoding/json"
	"fmt"
	"net/http"

	module "github.com/agile-work/srv-mdl-shared"
	"github.com/agile-work/srv-mdl-shared/models"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Update object data in the database
func Update(r *http.Request, object interface{}, scope, table string, condition builder.Builder) *module.Response {
	languageCode := r.Header.Get("Content-Language")
	models.TranslationFieldsRequestLanguageCode = languageCode
	response := GetResponse(r, object, scope)
	if response.Code != http.StatusOK {
		return response
	}

	getTranslationColumns := true
	if languageCode == "all" {
		getTranslationColumns = false
	}
	columns, translationColumns, bodyValues := getColumnsFromBody(r, object, getTranslationColumns)

	// TODO: change to db transaction to avoid updating one part without the translations
	if len(columns) > 0 {
		err := db.UpdateStruct(table, object, condition, columns...)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update", scope), err.Error()))

			return response
		}
	}

	if len(translationColumns) > 0 {
		statement := builder.Update(table)
		for _, col := range translationColumns {
			statement.JSON(col, languageCode)
			val, _ := json.Marshal(bodyValues[col])
			statement.Values(val)
		}
		statement.Where(condition)

		err := db.Exec(statement)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update translations", scope), err.Error()))

			return response
		}
	}

	response.Data = object

	return response
}
