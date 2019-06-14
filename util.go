package shared

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/agile-work/srv-mdl-shared/models"
)

func GetColumnsFromBody(body []byte, object interface{}) ([]string, map[string]string, error) {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &jsonMap); err != nil {
		return nil, nil, err
	}
	objectTranslationColumns := []string{}
	if models.TranslationFieldsRequestLanguageCode != "all" {
		objectTranslationColumns = getObjectTranslationColumns(object)
	}
	columns := []string{}
	translations := make(map[string]string)
	for k, v := range jsonMap {
		if k != "created_by" && k != "created_at" && k != "updated_by" && k != "updated_at" && !isValueInList(k, objectTranslationColumns) {
			columns = append(columns, k)
		} else if isValueInList(k, objectTranslationColumns) {
			translations[k] = v.(string)
		}
	}
	return columns, translations, nil
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
		if elementType.Field(i).Type == reflect.TypeOf(models.Translation{}) {
			translationColumns = append(translationColumns, elementType.Field(i).Tag.Get("sql"))
		}
	}
	return translationColumns
}

func isValueInList(value string, list []string) bool {
	if list == nil || len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
