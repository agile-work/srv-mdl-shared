package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	module "github.com/agile-work/srv-mdl-shared"
	shared "github.com/agile-work/srv-shared"
)

// LoadBodyToStruct load the request body to an object
func LoadBodyToStruct(r *http.Request, object interface{}) error {
	body, _ := GetBody(r)

	err := json.Unmarshal(body, &object)
	if err != nil {
		return err
	}
	return nil
}

// GetResponse load request body to object and creates a response
func GetResponse(r *http.Request, object interface{}, scope string) *module.Response {
	response := &module.Response{
		Code: http.StatusOK,
	}

	body, _ := GetBody(r)
	if len(body) > 0 {
		err := json.Unmarshal(body, object)
		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Errors = append(response.Errors, module.NewResponseError(shared.ErrorParsingRequest, fmt.Sprintf("%s unmarshal body", scope), err.Error()))

			return response
		}
	}

	SetSchemaAudit(r, object)

	return response
}

// SetSchemaAudit load user and time to audit fields
func SetSchemaAudit(r *http.Request, object interface{}) {
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
}

// GetBody get request body while maintaining the value in the request
func GetBody(r *http.Request) ([]byte, error) {
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// getColumnsFromBody get request body and return an string array with columns from the body
func getColumnsFromBody(r *http.Request, object interface{}) []string {
	body, _ := GetBody(r)
	jsonMap := make(map[string]interface{})
	json.Unmarshal(body, &jsonMap)
	columns := []string{}
	for k := range jsonMap {
		if k != "created_by" && k != "created_at" && k != "updated_by" && k != "updated_at" {
			columns = append(columns, k)
		}
	}

	elementValue := reflect.ValueOf(object).Elem()

	if r.Method == http.MethodPost {
		elementCreatedBy := elementValue.FieldByName("CreatedBy")
		elementCreatedAt := elementValue.FieldByName("CreatedAt")
		if elementCreatedBy.IsValid() {
			columns = append(columns, "created_by")
		}
		if elementCreatedAt.IsValid() {
			columns = append(columns, "created_at")
		}
	}

	elementUpdatedBy := elementValue.FieldByName("UpdatedBy")
	elementUpdatedAt := elementValue.FieldByName("UpdatedAt")
	if elementUpdatedBy.IsValid() {
		columns = append(columns, "updated_by")
	}
	if elementUpdatedAt.IsValid() {
		columns = append(columns, "updated_at")
	}

	return columns
}
