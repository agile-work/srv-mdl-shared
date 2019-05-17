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
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Update object data in the database
func Update(r *http.Request, object interface{}, scope, table string, condition builder.Builder) *module.Response {
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

	columns := getUpdateColumnsFromBody(body)

	userID := r.Header.Get("userID")
	now := time.Now()
	elementValue := reflect.ValueOf(object).Elem()
	elementUpdatedBy := elementValue.FieldByName("UpdatedBy")
	elementUpdatedAt := elementValue.FieldByName("UpdatedAt")
	elementUpdatedBy.SetString(userID)
	elementUpdatedAt.Set(reflect.ValueOf(now))

	err = db.UpdateStruct(table, object, condition, columns...)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update", scope), err.Error()))

		return response
	}

	translationColumns := GetTranslationLanguageCodeColumns(object, columns...)

	if len(translationColumns) > 0 {
		err = UpdateTranslationsFromStruct(table, r.Header.Get("Content-Language"), object, columns...)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorInsertingRecord, fmt.Sprintf("%s update translation", scope), err.Error()))

			return response
		}
	}

	response.Data = object

	return response
}

// getUpdateColumnsFromBody get request body and return an string array with columns from the body
func getUpdateColumnsFromBody(body []byte) []string {
	jsonMap := make(map[string]interface{})
	json.Unmarshal(body, &jsonMap)
	columns := []string{}
	for k := range jsonMap {
		if k != "created_by" && k != "created_at" && k != "updated_by" && k != "updated_at" {
			columns = append(columns, k)
		}
	}
	columns = append(columns, "updated_by")
	columns = append(columns, "updated_at")

	return columns
}
