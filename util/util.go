package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/util"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// LoadSQLOptionsFromURLQuery load sql options from url query
func LoadSQLOptionsFromURLQuery(query url.Values, opt *db.Options) {
	opt.Limit, _ = strconv.Atoi(query.Get("limit"))
	opt.Offset, _ = strconv.Atoi(query.Get("offset"))
}

// GetColumnsFromBody get a body and return an string array with columns from the body
func GetColumnsFromBody(body map[string]interface{}, object interface{}) ([]string, map[string]string) {
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

	return columns, translations
}

// GetBodyMap get request body while maintaining the value in the request
func GetBodyMap(r *http.Request) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	body, err := GetBody(r)
	if err != nil {
		return nil, customerror.New(http.StatusBadRequest, "GetBodyMap get body", err.Error())
	}

	if len(body) > 0 {
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return nil, customerror.New(http.StatusBadRequest, "GetBodyMap parse", err.Error())
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
			return nil, customerror.New(http.StatusBadRequest, "GetBody read", err.Error())
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
func SetSchemaAudit(isCreate bool, username string, object interface{}) {
	now := time.Now()
	elementValue := reflect.ValueOf(object).Elem()

	if isCreate {
		elementCreatedBy := elementValue.FieldByName("CreatedBy")
		elementCreatedAt := elementValue.FieldByName("CreatedAt")
		if elementCreatedBy.IsValid() {
			elementCreatedBy.SetString(username)
		}
		if elementCreatedAt.IsValid() {
			if elementCreatedAt.Kind() == reflect.Ptr {
				elementCreatedAt.Set(reflect.ValueOf(&now))
			} else {
				elementCreatedAt.Set(reflect.ValueOf(now))
			}
		}
	}

	elementUpdatedBy := elementValue.FieldByName("UpdatedBy")
	elementUpdatedAt := elementValue.FieldByName("UpdatedAt")
	if elementUpdatedBy.IsValid() {
		elementUpdatedBy.SetString(username)
	}
	if elementUpdatedAt.IsValid() {
		if elementUpdatedAt.Kind() == reflect.Ptr {
			elementUpdatedAt.Set(reflect.ValueOf(&now))
		} else {
			elementUpdatedAt.Set(reflect.ValueOf(now))
		}
	}
}

// Unique returns a slice with unique items
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// GetContentPrefix validate if content exists and return the prefix
func GetContentPrefix(code string) (string, error) {
	rows, err := db.Query(builder.Select("prefix", "is_module", "is_system").From(constants.TableCoreContents).Where(builder.Equal("code", code)))
	if err != nil {
		return "", err
	}

	content, err := db.MapScan(rows)
	if err != nil {
		return "", err
	}

	if len(content) <= 0 {
		return "", fmt.Errorf("invalid code")
	}

	prefix := content[0]["prefix"].(string)
	module := ""
	if content[0]["is_module"].(bool) {
		module = "mdl_"
	}
	if content[0]["is_system"].(bool) {
		prefix = "sys_" + module + prefix
	} else {
		prefix = "custom_" + prefix
	}

	return prefix, nil
}

// ValidateContent validate if the content exists in the database
func ValidateContent(code string) error {
	total, err := db.Count("code", constants.TableCoreContents, &db.Options{
		Conditions: builder.Equal("code", code),
	})
	if err != nil {
		return err
	}
	if total <= 0 {
		return fmt.Errorf("code not found")
	}
	return nil
}

// GetBodyUpdatableJSONColumns get all columns from body based on the struct fields
func GetBodyUpdatableJSONColumns(r *http.Request, isCreate bool, object interface{}, username, languageCode string) (map[string]interface{}, error) {
	bodyMap, err := GetBodyMap(r)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	elem := reflect.TypeOf(object).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Tag.Get("updatable") != "false" {
			col := strings.Split(field.Tag.Get("json"), ",")[0]
			if val, ok := bodyMap[col]; ok {
				if field.Type == reflect.TypeOf(translation.Translation{}) {
					path := col
					if languageCode != "all" {
						path = fmt.Sprintf("%s, %s", col, languageCode)
					}
					result[path] = val
				} else {
					result[col] = val
				}
			}
		}
		now := time.Now()
		if isCreate && field.Name == "CreatedBy" {
			result["created_by"] = username
		}
		if isCreate && field.Name == "CreatedAt" {
			result["created_at"] = now
		}
		if field.Name == "UpdatedBy" {
			result["updated_by"] = username
		}
		if field.Name == "UpdatedAt" {
			result["updated_at"] = now
		}
	}

	return result, nil
}

// DataToStruct convert a map to struct
func DataToStruct(data interface{}, object interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dataByte, object); err != nil {
		return err
	}
	return nil
}
