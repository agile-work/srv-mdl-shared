package db

// StructurePermission defines attirbutes to a structure
type StructurePermission struct {
	ID             string `json:"id" sql:"id" pk:"true"`
	Code           string `json:"code" sql:"code"`
	StructureType  string `json:"structure_type" sql:"structure_type"`   // field, widget, section...
	StructureClass string `json:"structure_class" sql:"structure_class"` //field: lookup_tree, text, date, number, money | widget: chart, table
	Permission     int    `json:"permission" sql:"permission"`
}

// GetUserAvailableStructures returns a list of available structures based on user permission
// passing StructureType = "" returns all types of structures
func GetUserAvailableStructures(UserID, SchemaCode, StructureType string) []StructurePermission {
	return nil
}
