package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Create object data in the database
func Create(r *http.Request, object interface{}, scope, table string) *module.Response {
	response := &module.Response{
		Code: http.StatusOK,
	}
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &object)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorParsingRequest, fmt.Sprintf("%s unmarshal body", scope), err.Error()))

		return response
	}

	userID := r.Header.Get("userID")
	now := time.Now()
	elementValue := reflect.ValueOf(object).Elem()
	elementCreatedBy := elementValue.FieldByName("CreatedBy")
	elementUpdatedBy := elementValue.FieldByName("UpdatedBy")
	elementCreatedAt := elementValue.FieldByName("CreatedAt")
	elementUpdatedAt := elementValue.FieldByName("UpdatedAt")
	elementCreatedBy.SetString(userID)
	elementUpdatedBy.SetString(userID)
	elementCreatedAt.Set(reflect.ValueOf(now))
	elementUpdatedAt.Set(reflect.ValueOf(now))

	id, err := db.InsertStruct(table, object)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s create", scope), err.Error()))

		return response
	}

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
