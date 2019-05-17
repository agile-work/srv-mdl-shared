package db

import (
	"net/http"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Load object data from the database
func Load(r *http.Request, object interface{}, scope, table string, conditions builder.Builder) *module.Response {
	response := &module.Response{
		Code: http.StatusOK,
	}

	filterColumns, _ := GetFilterColumns(r, object, table)

	if len(filterColumns) > 0 {
		newCondition := []builder.Builder{}
		if conditions != nil {
			newCondition = append(newCondition, conditions)
		}
		for column, value := range filterColumns {
			newCondition = append(newCondition, builder.Equal(column, value))
		}
		conditions = builder.And(newCondition...)
	}

	translationColumns := GetTranslationLanguageCodeColumns(object)

	if len(translationColumns) > 0 {
		newCondition := []builder.Builder{}
		if conditions != nil {
			newCondition = append(newCondition, conditions)
		}
		for _, translationColumn := range translationColumns {
			newCondition = append(newCondition, builder.Equal(translationColumn, r.Header.Get("Content-Language")))
		}
		conditions = builder.And(newCondition...)
	}

	err := db.LoadStruct(table, object, conditions)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorLoadingData, scope, err.Error()))

		return response
	}

	response.Data = object

	return response
}
