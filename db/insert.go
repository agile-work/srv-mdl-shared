package db

import (
	"fmt"
	"net/http"
	"reflect"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Create object data in the database
func Create(r *http.Request, object interface{}, scope, table string) *module.Response {
	response := GetResponse(r, object, scope)
	if response.Code != http.StatusOK {
		return response
	}

	columns := getColumnsFromBody(r, object)

	id, err := db.InsertStruct(table, object, columns...)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s create", scope), err.Error()))

		return response
	}

	elementValue := reflect.ValueOf(object).Elem()
	elementID := elementValue.FieldByName("ID")
	elementID.SetString(id)

	translationColumns := GetTranslationLanguageCodeColumns(object)

	if len(translationColumns) > 0 {
		err = CreateTranslationsFromStruct(table, r.Header.Get("Content-Language"), object)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s create translation", scope), err.Error()))

			return response
		}
	}

	response.Data = object

	return response
}
