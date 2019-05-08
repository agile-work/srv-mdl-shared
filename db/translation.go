package db

import (
	"fmt"
	"reflect"

	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Translation defines the struct of this object
type Translation struct {
	ID             string `json:"id" sql:"id" pk:"true"`
	StructureID    string `json:"structure_id" sql:"structure_id" fk:"true"`
	StructureType  string `json:"structure_type" sql:"structure_type"`
	StructureField string `json:"structure_field" sql:"structure_field"`
	Value          string `json:"value" sql:"value"`
	LanguageCode   string `json:"language_code" sql:"language_code"`
	Replicated     bool   `json:"replicated" sql:"replicated"`
}

// GetTranslationLanguageCodeColumns return translation columns from the object
func GetTranslationLanguageCodeColumns(object interface{}, columns ...string) []string {
	translationColumns := []string{}
	elementType := reflect.TypeOf(object).Elem()

	if elementType.Kind() == reflect.Slice {
		elementType = elementType.Elem()
	}

	for i := 0; i < elementType.NumField(); i++ {
		elementField := elementType.Field(i)
		if elementField.Tag.Get("table") == shared.TableCoreTranslations {
			jsonColumn := elementField.Tag.Get("json")
			translationTableAlias := elementField.Tag.Get("alias")
			if len(columns) > 0 {
				for _, column := range columns {
					if column == jsonColumn {
						translationColumns = append(translationColumns, fmt.Sprintf("%s.language_code", translationTableAlias))
					}
				}
			} else {
				translationColumns = append(translationColumns, fmt.Sprintf("%s.language_code", translationTableAlias))
			}
		}
	}

	return translationColumns
}

// CreateTranslationsFromStruct saves translations from struct to the database
func CreateTranslationsFromStruct(structureType, languageCode string, object interface{}) error {
	objectType := reflect.TypeOf(object).Elem()
	objectValue := reflect.ValueOf(object).Elem()

	translations := []Translation{}
	for i := 0; i < objectType.NumField(); i++ {
		if objectType.Field(i).Tag.Get("table") == shared.TableCoreTranslations {
			structureID := objectValue.FieldByName("ID").Interface().(string)
			structureField := objectType.Field(i).Tag.Get("json")
			value := objectValue.Field(i).Interface().(string)
			translation := Translation{
				StructureID:    structureID,
				StructureField: structureField,
				StructureType:  structureType,
				Value:          value,
				LanguageCode:   languageCode,
			}
			translations = append(translations, translation)
		}
	}

	_, err := db.InsertStruct(shared.TableCoreTranslations, translations)
	return err
}

// UpdateTranslationsFromStruct updates translations from struct to the database
func UpdateTranslationsFromStruct(structureType, languageCode string, object interface{}, columns ...string) error {
	objectType := reflect.TypeOf(object).Elem()
	objectValue := reflect.ValueOf(object).Elem()

	for i := 0; i < objectType.NumField(); i++ {
		objectField := objectType.Field(i)
		if objectField.Tag.Get("table") == shared.TableCoreTranslations {
			for _, column := range columns {
				if column == objectField.Tag.Get("json") {
					structureID := objectValue.FieldByName("ID").Interface().(string)
					structureField := objectField.Tag.Get("json")
					value := objectValue.Field(i).Interface().(string)
					translation := Translation{
						StructureID:    structureID,
						StructureField: structureField,
						StructureType:  structureType,
						Value:          value,
						LanguageCode:   languageCode,
					}

					structureIDColumn := fmt.Sprintf("%s.structure_id", shared.TableCoreTranslations)
					structureFieldColumn := fmt.Sprintf("%s.structure_field", shared.TableCoreTranslations)
					languageCodeColumn := fmt.Sprintf("%s.language_code", shared.TableCoreTranslations)
					condition := builder.And(
						builder.Equal(structureIDColumn, structureID),
						builder.Equal(structureFieldColumn, structureField),
						builder.Equal(languageCodeColumn, languageCode),
					)

					err := db.UpdateStruct(
						shared.TableCoreTranslations, &translation, condition,
						"structure_id", "structure_field", "structure_type", "value", "language_code",
					)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}
