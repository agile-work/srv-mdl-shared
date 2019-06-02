package models

const (
	// TODO: Get default system language from a configuration file
	systemDefaultLanguage string = "pt-br"
)

// Translation represents a translation json object
type Translation map[string]string

// // UnmarshalJSON custom
// func (t Translation) UnmarshalJSON(data []byte) error {
// 	mapJSON := make(map[string]string)
// 	mapJSON["undefined"] = string(data)
// 	t = Translation(mapJSON)
// 	return nil
// }

// String returns the translation value for a language code
func (t Translation) String(code string) string {
	if val, ok := t[code]; ok {
		return val
	}
	return ""
}

// GetAvailable always return a value for a translation field
func (t Translation) GetAvailable(code string) string {
	val := t.String(code)
	if val == "" {
		val := t.String(systemDefaultLanguage)
		if val == "" {
			for _, t := range t {
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
	for k, v := range t {
		if k == "undefined" {
			t[code] = v
			delete(t, k)
		}
	}
}
