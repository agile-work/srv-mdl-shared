package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/agile-work/srv-shared/util"

	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// LoadSQLOptionsFromURLQuery load sql options from url query
func LoadSQLOptionsFromURLQuery(query url.Values, opt *db.Options) {
	opt.Limit, _ = strconv.Atoi(query.Get("limit"))
	opt.Offset, _ = strconv.Atoi(query.Get("offset"))
}

// GetColumnsFromBody get a body and return an string array with columns from the body
func GetColumnsFromBody(body map[string]interface{}, object interface{}) ([]string, map[string]string, error) {
	objectTranslationColumns := []string{}
	if translation.FieldsRequestLanguageCode != "all" {
		objectTranslationColumns = getObjectTranslationColumns(object)
	}
	columns := []string{}
	translations := make(map[string]string)
	for k, v := range body {
		if k != "created_by" && k != "created_at" && k != "updated_by" && k != "updated_at" && !util.Contains(objectTranslationColumns, k) {
			columns = append(columns, k)
		} else if util.Contains(objectTranslationColumns, k) {
			translations[k] = v.(string)
		}
	}
	columns = append(columns, "updated_by")
	columns = append(columns, "updated_at")

	return columns, translations, nil
}

// GetBodyMap get request body while maintaining the value in the request
func GetBodyMap(r *http.Request) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	body, err := GetBody(r)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return nil, err
		}
	}

	return jsonMap, nil
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

// getObjectTranslationColumns return an array with all translation columns from an object
func getObjectTranslationColumns(object interface{}) []string {
	translationColumns := []string{}
	elementType := reflect.TypeOf(object).Elem()
	for i := 0; i < elementType.NumField(); i++ {
		if elementType.Field(i).Type == reflect.TypeOf(translation.Translation{}) {
			translationColumns = append(translationColumns, elementType.Field(i).Tag.Get("sql"))
		}
	}
	return translationColumns
}

// GetBodyColumns return all columns from body
func GetBodyColumns(body map[string]interface{}) []string {
	columns := []string{}
	for k := range body {
		columns = append(columns, k)
	}
	return columns
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
