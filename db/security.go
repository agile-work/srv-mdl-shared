package db

import (
	"fmt"

	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	sql "github.com/agile-work/srv-shared/sql-builder/db"
)

// StructurePermission defines attirbutes to a structure
type StructurePermission struct {
	ID             string `json:"id" sql:"id" pk:"true"`
	StructureCode  string `json:"structure_code" sql:"structure_code"`
	StructureType  string `json:"structure_type" sql:"structure_type"`   // field, widget, section...
	StructureClass string `json:"structure_class" sql:"structure_class"` //field: lookup_tree, text, date, number, money | widget: chart, table
	Permission     int    `json:"permission" sql:"permission"`
}

// GetUserAvailableFields returns a list of available structures based on user permission
// passing StructureType = "" returns all types of structures
func GetUserAvailableFields(userID, schemaCode, structureType string) ([]StructurePermission, error) {
	permissions := []StructurePermission{}
	userIDColumn := fmt.Sprintf("%s.user_id", shared.ViewCoreUserStructurePermissions)
	schemaCodeColumn := fmt.Sprintf("%s.schema_code", shared.ViewCoreUserStructurePermissions)
	structureTypeColumn := fmt.Sprintf("%s.structure_type", shared.ViewCoreUserStructurePermissions)
	condition := builder.And(
		builder.Equal(userIDColumn, userID),
		builder.Equal(schemaCodeColumn, schemaCode),
		builder.Equal(structureTypeColumn, structureType),
	)

	err := sql.LoadStruct(shared.ViewCoreUserStructurePermissions, &permissions, condition)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}
