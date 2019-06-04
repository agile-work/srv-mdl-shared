package db

import (
	"fmt"
	"net/http"
	"reflect"

	module "github.com/agile-work/srv-mdl-shared"
	"github.com/agile-work/srv-mdl-shared/models"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Create object data in the database
func Create(r *http.Request, object interface{}, scope, table string) *module.Response {
	models.TranslationFieldsRequestLanguageCode = r.Header.Get("Content-Language")
	response := GetResponse(r, object, scope)
	if response.Code != http.StatusOK {
		return response
	}

	columns, _, _ := getColumnsFromBody(r, object, false)
	models.TranslationFieldsRequestLanguageCode = "all"

	id, err := db.InsertStruct(table, object, columns...)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s create", scope), err.Error()))

		return response
	}

	elementValue := reflect.ValueOf(object).Elem()
	elementID := elementValue.FieldByName("ID")
	elementID.SetString(id)

	response.Data = object

	return response
}
