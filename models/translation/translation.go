package translation

import (
	"encoding/json"
)

// FieldsRequestLanguageCode should be set when passing a translation fileds in payload with only a string
// If not defined the value "undefined" will be used
var FieldsRequestLanguageCode = "undefined"

// SystemDefaultLanguageCode defines the system language code
var SystemDefaultLanguageCode = "pt-br"

// Translation represents a translation json object
type Translation struct {
	Language map[string]string `json:"languages"`
}

// UnmarshalJSON custom unmarshal function to deal with multiple translations or just a string
func (t *Translation) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &t.Language)
	if err != nil {
		t.Language = make(map[string]string)
		val := string(data)
		t.Language[FieldsRequestLanguageCode] = val[1 : len(val)-1]
	}
	return nil
}

// MarshalJSON custom marshal function to return only the map of available languages
func (t Translation) MarshalJSON() ([]byte, error) {
	if FieldsRequestLanguageCode != "undefined" && FieldsRequestLanguageCode != "all" {
		return json.Marshal(t.GetAvailable(FieldsRequestLanguageCode))
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
