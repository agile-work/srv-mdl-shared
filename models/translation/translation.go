package translation

import (
	"encoding/json"
	"reflect"
)

// FieldsRequestLanguageCode should be set when passing a translation fileds in payload with only a string
// If not defined the value "undefined" will be used
// TODO: refactoring delete
var FieldsRequestLanguageCode = "undefined"

// SystemDefaultLanguageCode defines the system language code
var SystemDefaultLanguageCode = "pt-br"

// Translation represents a translation json object
type Translation struct {
	Language            map[string]string `json:"languages"`
	RequestLanguageCode string
}

// UnmarshalJSON custom unmarshal function to deal with multiple translations or just a string
func (t *Translation) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &t.Language)
	if err != nil {
		t.Language = make(map[string]string)
		val := string(data)
		if t.RequestLanguageCode == "" {
			t.RequestLanguageCode = FieldsRequestLanguageCode
		}
		t.Language[t.RequestLanguageCode] = val[1 : len(val)-1]
	}
	return nil
}

// MarshalJSON custom marshal function to return only the map of available languages
func (t Translation) MarshalJSON() ([]byte, error) {
	if t.RequestLanguageCode == "" {
		t.RequestLanguageCode = FieldsRequestLanguageCode
	}
	if t.RequestLanguageCode != "" && t.RequestLanguageCode != "all" {
		return json.Marshal(t.GetAvailable(t.RequestLanguageCode))
	}
	return json.Marshal(t.Language)
}

// String returns the translation value for a language code
func (t Translation) String(code string) string {
	if val, ok := t.Language[code]; ok {
		return val
	}
	return ""
}

// GetAvailable always return a value for a translation field
func (t Translation) GetAvailable(code string) string {
	val := t.String(code)
	if val == "" {
		val := t.String(SystemDefaultLanguageCode)
		if val == "" {
			for _, t := range t.Language {
				val = t
				return val
			}
		}
		return val
	}
	return val
}

// Parse put the payload in the correct format
func (t Translation) Parse(code string) {
	for k, v := range t.Language {
		if k == "undefined" {
			t.Language[code] = v
			delete(t.Language, k)
		}
	}
}

// SetStructTranslationsLanguage check struct translations fields and set the language code
func SetStructTranslationsLanguage(object interface{}, languageCode string) {
	inspectStruct(reflect.ValueOf(object), languageCode)
}

// SetSliceTranslationsLanguage check struct translations fields and set the language code
func SetSliceTranslationsLanguage(slice interface{}, languageCode string) {
	s := reflect.ValueOf(slice).Elem()
	for i := 0; i < s.Len(); i++ {
		inspectStruct(s.Index(i), languageCode)
	}
}

func inspectStruct(val reflect.Value, languageCode string) {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}
		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()
		}
		if valueField.Kind() == reflect.Struct {
			if valueField.Type() == reflect.TypeOf(Translation{}) {
				if isZero(valueField) {
					valueField.Set(reflect.ValueOf(Translation{RequestLanguageCode: languageCode}))
				} else {
					valueField.FieldByName("RequestLanguageCode").SetString(languageCode)
				}
			} else {
				inspectStruct(valueField, languageCode)
			}

		}
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				z = z && isZero(v.Field(i))
			}
		}
		return z
	case reflect.Ptr:
		return isZero(reflect.Indirect(v))
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	result := v.Interface() == z.Interface()

	return result
}
