package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/agile-work/srv-shared/sql-builder/db"

	"github.com/agile-work/srv-shared/sql-builder/builder"

	shared "github.com/agile-work/srv-shared"
	jwt "github.com/dgrijalva/jwt-go"
)

// User defines the struct of this object
type User struct {
	ID                string            `json:"id" sql:"id" pk:"true"`
	Username          string            `json:"username" sql:"username"`
	FirstName         string            `json:"first_name" sql:"first_name"`
	LastName          string            `json:"last_name" sql:"last_name"`
	Email             string            `json:"email" sql:"email"`
	Password          string            `json:"password" sql:"password"`
	LanguageCode      string            `json:"language_code" sql:"language_code"`
	ReceiveEmails     string            `json:"receive_emails" sql:"receive_emails"`
	Security          security          `json:"security" sql:"security" field:"jsonb"`
	SecurityInstances securityInstances `json:"security_instances" sql:"security_instances" field:"jsonb"`
	Active            bool              `json:"active" sql:"active"`
	CreatedBy         string            `json:"created_by" sql:"created_by"`
	CreatedAt         time.Time         `json:"created_at" sql:"created_at"`
	UpdatedBy         string            `json:"updated_by" sql:"updated_by"`
	UpdatedAt         time.Time         `json:"updated_at" sql:"updated_at"`
}

// GetUserSelectableFields returns an array with the user fields
func GetUserSelectableFields() []string {
	return []string{
		shared.TableCoreUsers + ".id as user_id",
		shared.TableCoreUsers + ".username",
		shared.TableCoreUsers + ".first_name",
		shared.TableCoreUsers + ".last_name",
		shared.TableCoreUsers + ".email",
		shared.TableCoreUsers + ".language_code",
		shared.TableCoreUsers + ".receive_emails",
	}
}

type securityInstanceDefinition struct {
	PermissionScope string          `json:"permission_scope"`
	Permission      string          `json:"permission"`
	View            map[string]bool `json:"view"`
	Edit            map[string]bool `json:"edit"`
}

type securityInstance struct {
	Instance map[string]securityInstanceDefinition `json:"instance"`
}

type securityInstances struct {
	Schema map[string]securityInstance `json:"schema"`
}

type securityDefinition struct {
	PermissionInstance  string            `json:"permission_instance"`
	SecurityFields      map[string]string `json:"security_fields"`
	UserTreesSecurity   map[string]string `json:"user_trees_security"`
	PermissionStructure string            `json:"permission_structure"`
	QueryDataColumns    string            `json:"query_data_columns"`
	View                map[string]bool   `json:"view"`
	Edit                map[string]bool   `json:"edit"`
}

type security struct {
	Schema map[string]securityDefinition `json:"schema"`
}

// ViewUserAllPermissions defines the struct of this object
type ViewUserAllPermissions struct {
	UserID         string `json:"user_id" sql:"user_id"`
	SchemaID       string `json:"schema_id" sql:"schema_id"`
	SchemaCode     string `json:"schema_code" sql:"schema_code"`
	SchemaName     string `json:"schema_name" sql:"schema_name"`
	StructureID    string `json:"structure_id" sql:"structure_id"`
	StructureCode  string `json:"structure_code" sql:"structure_code"`
	StructureType  string `json:"structure_type" sql:"structure_type"`
	StructureClass string `json:"structure_class" sql:"structure_class"`
	StructureName  string `json:"structure_name" sql:"structure_name"`
	LanguageCode   string `json:"language_code" sql:"language_code"`
	PermissionType int    `json:"permission_type" sql:"permission_type"`
	Scope          string `json:"scope" sql:"scope"`
}

// ViewGroupUser defines the struct of this object
type ViewGroupUser struct {
	ID            string    `json:"id" sql:"id" pk:"true"`
	GroupID       string    `json:"group_id" sql:"group_id" fk:"true"`
	Username      string    `json:"username" sql:"username"`
	FirstName     string    `json:"first_name" sql:"first_name"`
	LastName      string    `json:"last_name" sql:"last_name"`
	Email         string    `json:"email" sql:"email"`
	Password      string    `json:"password" sql:"password"`
	LanguageCode  string    `json:"language_code" sql:"language_code"`
	ReceiveEmails string    `json:"receive_emails" sql:"receive_emails"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// UserCustomClaims used to parse token payload
type UserCustomClaims struct {
	User User
	jwt.StandardClaims
}

// InitSecurityQuery return the columns to make a security query
func (u *User) InitSecurityQuery(schemaCode string) (*builder.Statement, error) {
	schemaTable := fmt.Sprintf("%s%s AS sch", shared.InstancesTablePrefix, schemaCode)
	securitySchema := u.Security.Schema[schemaCode]
	columnsJSON := strings.Split(securitySchema.QueryDataColumns, ",")
	columns := []string{}

	columns = append(columns, "sch.id")

	if columnsJSON[0] == "*" {
		columnsJSON = []string{}
		statement := builder.Select(
			"fld.code",
		).From(
			fmt.Sprintf("%s AS %s", shared.TableCoreSchemaFields, "fld"),
		).Join(
			fmt.Sprintf("%s AS %s", shared.TableCoreSchemas, "sch"),
			"sch.id = fld.schema_id",
		).Where(
			builder.Equal("sch.code", schemaCode),
		)

		rows, err := db.Query(statement)
		if err != nil {
			return nil, err
		}

		results, err := db.MapScan(rows)
		if err != nil {
			return nil, err
		}

		for _, row := range results {
			columnsJSON = append(columnsJSON, row["code"].(string))
		}
	}

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree := range securitySchema.SecurityFields {
			columns = append(columns, fmt.Sprintf("unit_%s.path AS unit_%s", schemaColumnTree, schemaColumnTree))
		}
	}

	return builder.Select(columns...).JSON("data", columnsJSON...).From(schemaTable), nil
}

// LoadSecurityTreeJoins return the tree joins to make a security query
func (u *User) LoadSecurityTreeJoins(schemaCode string, statementSchema *builder.Statement) {
	securitySchema := u.Security.Schema[schemaCode]

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree, tree := range securitySchema.SecurityFields {
			statementSchema.LeftJoin(
				fmt.Sprintf("%s AS tree_%s", shared.TableCoreTrees, schemaColumnTree),
				fmt.Sprintf("tree_%s.code = '%s'", schemaColumnTree, tree),
			)
			statementSchema.LeftJoin(
				fmt.Sprintf("%s AS unit_%s", shared.TableCoreTreeUnits, schemaColumnTree),
				fmt.Sprintf(
					"unit_%s.tree_id = tree_%s.id AND unit_%s.code = sch.data->>'%s'",
					schemaColumnTree,
					schemaColumnTree,
					schemaColumnTree,
					schemaColumnTree,
				),
			)
		}
	}
}

// LoadSecurityConditions return the where condition to make a security query
func (u *User) LoadSecurityConditions(schemaCode string, statementSchema *builder.Statement) {
	securitySchema := u.Security.Schema[schemaCode]
	securityInstanceSchema := u.SecurityInstances.Schema[schemaCode]
	conditions := []builder.Builder{}

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree, tree := range securitySchema.SecurityFields {
			if treePath, ok := securitySchema.UserTreesSecurity[tree]; ok {
				conditions = append(
					conditions,
					builder.Raw(
						fmt.Sprintf(
							"unit_%s.path ~ '%s'",
							schemaColumnTree,
							treePath,
						),
					),
				)
			}
		}

		instances := []string{}
		for instance := range securityInstanceSchema.Instance {
			instances = append(instances, fmt.Sprintf("'%s'", instance))
		}

		if len(instances) > 0 {
			conditions = append(
				conditions,
				builder.Raw(
					fmt.Sprintf(
						"sch.id IN (%s)",
						strings.Join(instances, ", "),
					),
				),
			)
		}

		statementSchema.Where(builder.Or(conditions...))
	}
}
