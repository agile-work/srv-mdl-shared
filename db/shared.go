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
)

// GetResponse load request body to object and creates a response
func GetResponse(r *http.Request, object interface{}, scope string) *module.Response {
	response := &module.Response{
		Code: http.StatusOK,
	}

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		err := json.Unmarshal(body, &object)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorParsingRequest, fmt.Sprintf("%s unmarshal body", scope), err.Error()))

			return response
		}
	}

	userID := r.Header.Get("userID")
	now := time.Now()
	elementValue := reflect.ValueOf(object).Elem()

	if r.Method == http.MethodPost {
		elementCreatedBy := elementValue.FieldByName("CreatedBy")
		elementCreatedAt := elementValue.FieldByName("CreatedAt")
		if elementCreatedBy.IsValid() {
			elementCreatedBy.SetString(userID)
		}
		if elementCreatedAt.IsValid() {
			elementCreatedAt.Set(reflect.ValueOf(now))
		}
	}

	elementUpdatedBy := elementValue.FieldByName("UpdatedBy")
	elementUpdatedAt := elementValue.FieldByName("UpdatedAt")
	if elementUpdatedBy.IsValid() {
		elementUpdatedBy.SetString(userID)
	}
	if elementUpdatedAt.IsValid() {
		elementUpdatedAt.Set(reflect.ValueOf(now))
	}

	return response
}
