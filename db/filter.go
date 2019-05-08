package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	shared "github.com/agile-work/srv-shared"
)

// GetFilterColumns return translation columns from the object
func GetFilterColumns(r *http.Request, object interface{}, table string) (map[string]interface{}, error) {
	query := r.URL.Query()
	jsonFilters := query.Get("filter")
	filterColumns := make(map[string]interface{})

	if jsonFilters != "" {
		data := []byte(jsonFilters)
		filterMap := make(map[string]interface{})
		err := json.Unmarshal(data, &filterMap)

		if err != nil {
			return nil, err
		}

		elementType := reflect.TypeOf(object).Elem()

		if elementType.Kind() == reflect.Slice {
			elementType = elementType.Elem()
		}

		for filter, value := range filterMap {
			column := fmt.Sprintf("%s.%s", table, filter)
			for i := 0; i < elementType.NumField(); i++ {
				elementField := elementType.Field(i)
				if filter == elementField.Tag.Get("json") && elementField.Tag.Get("table") == shared.TableCoreTranslations {
					column = fmt.Sprintf("%s.%s", elementField.Tag.Get("alias"), "value")
					break
				}
			}
			filterColumns[column] = value
		}
	}

	return filterColumns, nil
}
