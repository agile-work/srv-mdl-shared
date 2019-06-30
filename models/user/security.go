package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agile-work/srv-shared/rdb"
	"github.com/agile-work/srv-shared/util"
	"github.com/tidwall/gjson"

	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

type securityInstanceDefinition struct {
	PermissionScope string          `json:"permission_scope"`
	Permission      string          `json:"permission"`
	Edit            map[string]bool `json:"edit"`
	View            map[string]bool `json:"view"`
}

type securityInstance struct {
	Instance map[string]securityInstanceDefinition `json:"instance"`
}

type securityInstances struct {
	Schema map[string]securityInstance `json:"schema"`
}

type securityTree struct {
	Tree       string          `json:"tree"`
	TreeUnit   string          `json:"tree_unit"`
	Permission string          `json:"permission"`
	Edit       map[string]bool `json:"edit"`
	View       map[string]bool `json:"view"`
}

type securityDefinition struct {
	PermissionInstance  string            `json:"permission_instance"`
	SecurityFields      map[string]string `json:"security_fields"`
	UserTreesSecurity   map[string]string `json:"user_trees_security"`
	PermissionStructure string            `json:"permission_structure"`
	QueryDataColumns    string            `json:"query_data_columns"`
	Edit                map[string]bool   `json:"edit"`
	View                map[string]bool   `json:"view"`
	Trees               []securityTree    `json:"trees"`
}

type security struct {
	Schema map[string]securityDefinition `json:"schema"`
}

// GetSecurityInstances return the initial statement to make a security query
func (u *User) GetSecurityInstances(schemaCode string, opt *db.Options, subQuery *builder.Statement, securityFields map[string]map[string]string) ([]map[string]interface{}, error) {
	securitySchema := u.Security.Schema[schemaCode]
	securityInstanceSchema := u.SecurityInstances.Schema[schemaCode]

	statement, err := getSecurityStatement(securitySchema, securityInstanceSchema, schemaCode, opt, subQuery)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}

	results, err := u.applySecurity(securitySchema, securityInstanceSchema, schemaCode, rows, securityFields)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func getSecurityStatement(securitySchema securityDefinition, securityInstanceSchema securityInstance, schemaCode string, opt *db.Options, subQuery *builder.Statement) (*builder.Statement, error) {
	schemaTable := fmt.Sprintf("%s%s AS sch", constants.InstancesTablePrefix, schemaCode)
	columns := []string{}
	statement := &builder.Statement{}

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree := range securitySchema.SecurityFields {
			columns = append(columns, fmt.Sprintf("unit_%s.path::TEXT AS unit_%s", schemaColumnTree, schemaColumnTree))
		}
	}

	if subQuery != nil {
		columns = append(columns, "sub.*")
		statement = builder.Select(columns...).From(schemaTable).JoinSubQuery(
			"sub", subQuery, builder.Raw("sub.id = sch.id"),
		)
	} else {
		columns = append(columns, "sch.id")
		columnsJSON := strings.Split(securitySchema.QueryDataColumns, ",")
		if columnsJSON[0] == "*" {
			columnsJSON = []string{}
			// TODO: Tratar o erro no redis
			rdbSchemaDefFields, _ := rdb.Get("schema:def:contract:fields")

			if rdbSchemaDefFields != "" {
				fields := gjson.Get(rdbSchemaDefFields, "#.code")
				for _, field := range fields.Array() {
					columnsJSON = append(columnsJSON, field.String())
				}
			} else {
				rows, err := db.Query(
					builder.Select(
						"fld.*",
					).From(
						fmt.Sprintf("%s AS %s", constants.TableCoreSchemaFields, "fld"),
					).Join(
						fmt.Sprintf("%s AS %s", constants.TableCoreSchemas, "sch"),
						"sch.code = fld.schema_code",
					).Where(
						builder.And(
							builder.Equal("sch.code", schemaCode),
							builder.Equal("fld.active", true),
						),
					),
				)
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

				jsonBytes, err := json.Marshal(results)
				if err != nil {
					return nil, err
				}

				if err := rdb.Set("schema:def:contract:fields", string(jsonBytes), 0); err != nil {
					return nil, err
				}
			}
		}
		statement = builder.Select(columns...).JSON("data", columnsJSON...).From(schemaTable)
	}

	loadSecurityTreeJoins(securitySchema, statement)
	loadSecurityConditions(securitySchema, securityInstanceSchema, statement)

	if opt.Conditions != nil {
		statement.Where(opt.Conditions)
	}

	if opt.OrderBy != nil {
		statement.OrderBy(opt.OrderBy...)
	}

	statement.Limit(opt.Limit)
	statement.Offset(opt.Offset)

	return statement, nil
}

// loadSecurityTreeJoins put the tree joins on statement to make a security query
func loadSecurityTreeJoins(securitySchema securityDefinition, statementSchema *builder.Statement) {
	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree, tree := range securitySchema.SecurityFields {
			statementSchema.LeftJoin(
				fmt.Sprintf("%s AS tree_%s", constants.TableCoreTrees, schemaColumnTree),
				fmt.Sprintf("tree_%s.code = '%s'", schemaColumnTree, tree),
			)
			statementSchema.LeftJoin(
				fmt.Sprintf("%s AS unit_%s", constants.TableCoreTreeUnits, schemaColumnTree),
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

// loadSecurityConditions put the where condition on statement to make a security query
func loadSecurityConditions(securitySchema securityDefinition, securityInstanceSchema securityInstance, statementSchema *builder.Statement) {
	conditions := []builder.Builder{}

	if securitySchema.PermissionInstance == "custom" {
		// TODO: Usar o SecurityFields em um cache com REDIS no lugar da informação no security do usuário
		for schemaColumnTree, tree := range securitySchema.SecurityFields {
			if treePath, ok := securitySchema.UserTreesSecurity[tree]; ok {
				conditions = append(conditions, builder.Raw(fmt.Sprintf("unit_%s.path ~ '%s'", schemaColumnTree, treePath)))
			}
		}

		instances := []string{}
		for instance := range securityInstanceSchema.Instance {
			instances = append(instances, fmt.Sprintf("'%s'", instance))
		}

		if len(instances) > 0 {
			conditions = append(conditions, builder.Raw(fmt.Sprintf("sch.id IN (%s)", strings.Join(instances, ", "))))
		}

		statementSchema.Where(builder.Or(conditions...))
	}
}

// applySecurity checks the security in the instance columns and clears if user do not have permission
func (u *User) applySecurity(securitySchema securityDefinition, securityInstanceSchema securityInstance, schemaCode string, rows *sql.Rows, securityFields map[string]map[string]string, columns ...string) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	schemaFields := map[string]string{}
	requiredFields := getRequiredFields(securitySchema.SecurityFields)
	if len(securityFields) > 0 {
		schemaFields = securityFields[schemaCode]
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		allFields := map[string]bool{}
		mapJSON, err := util.RowToMap(rows)
		if err != nil {
			return nil, err
		}
		for _, column := range cols {
			if requiredFields[column] {
				if column != "id" {
					allFields[column] = false
				}
				continue
			}
			if len(columns) > 0 && !util.Contains(columns, column) {
				allFields[column] = false
				continue
			}
			if securitySchema.PermissionStructure == "custom" {
				if len(securityFields) == 0 {
					schemaFields = map[string]string{}
					schemaFields[column] = column
				}
				if securityColumn, ok := schemaFields[column]; ok {
					allFields[column] = false
					allFields = applySecurityGroup(securitySchema, securityInstanceSchema, column, securityColumn, mapJSON, allFields)
					allFields = applySecurityTree(securitySchema, column, securityColumn, mapJSON, allFields)
					allFields = applySecurityInstance(securityInstanceSchema, column, securityColumn, mapJSON, allFields)
				} else {
					allFields[column] = true
					continue
				}
			}
		}
		results = append(results, getSecurityDataFields(mapJSON, allFields))
	}
	rows.Close()

	return results, nil
}

func getRequiredFields(securityFields map[string]string) map[string]bool {
	requiredFields := map[string]bool{}
	requiredFields["id"] = true
	// TODO: Usar o SecurityFields em um cache com REDIS no lugar da informação no security do usuário
	for schemaColumnTree := range securityFields {
		// TODO: Usar o requiredFields em um cache com REDIS
		requiredFields[schemaColumnTree] = true
		requiredFields["unit_"+schemaColumnTree] = true
	}
	return requiredFields
}

func applySecurityGroup(securitySchema securityDefinition, securityInstanceSchema securityInstance, column, securityColumn string, data map[string]interface{}, fields map[string]bool) map[string]bool {
	if instanceOptions, ok := securityInstanceSchema.Instance[data["id"].(string)]; !ok || (ok && instanceOptions.PermissionScope != "replace") {
		if securitySchema.PermissionStructure == "custom" {
			if rule, ok := securitySchema.Edit[securityColumn]; ok {
				fields[column] = rule
			}
			if rule, ok := securitySchema.View[securityColumn]; ok {
				fields[column] = rule
			}

		} else {
			fields[column] = true
		}
	}
	return fields
}

func applySecurityTree(securitySchema securityDefinition, column, securityColumn string, data map[string]interface{}, fields map[string]bool) map[string]bool {
LoopTrees:
	for _, securitySchemaTree := range securitySchema.Trees {
		for schemaColumnTree, tree := range securitySchema.SecurityFields {
			if tree == securitySchemaTree.Tree {
				securityUnitPath := securitySchemaTree.TreeUnit
				userUnit := strings.Replace(
					strings.Replace(
						securitySchema.UserTreesSecurity[tree], "*", "", -1,
					), ".", "", -1,
				)
				instanceSecurityUnitPath := data[("unit_" + schemaColumnTree)].(string)
				if (securityUnitPath[len(securityUnitPath)-1:] == "*" &&
					(instanceSecurityUnitPath == userUnit ||
						strings.HasPrefix(instanceSecurityUnitPath, userUnit+".") ||
						strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit) ||
						strings.Contains(instanceSecurityUnitPath, "."+userUnit+"."))) ||
					(securityUnitPath[len(securityUnitPath)-1:] != "*" &&
						(instanceSecurityUnitPath == userUnit ||
							strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit))) {
					if securitySchemaTree.Permission != "custom" {
						fields[column] = true
					}
					if rule, ok := securitySchemaTree.Edit[securityColumn]; ok {
						fields[column] = rule
					}
					if rule, ok := securitySchemaTree.View[securityColumn]; ok {
						fields[column] = rule
					}
					break LoopTrees
				}
			}
		}
	}
	return fields
}

func applySecurityInstance(securityInstanceSchema securityInstance, column, securityColumn string, data map[string]interface{}, fields map[string]bool) map[string]bool {
	if instanceOptions, ok := securityInstanceSchema.Instance[data["id"].(string)]; ok {
		if instanceOptions.PermissionScope == "replace" {
			fields[column] = false
		}
		if instanceOptions.Permission != "custom" {
			fields[column] = true
		}
		if rule, ok := instanceOptions.Edit[securityColumn]; ok {
			fields[column] = rule
		}
		if rule, ok := instanceOptions.View[securityColumn]; ok {
			fields[column] = rule
		}
	}
	return fields
}

func getSecurityDataFields(data map[string]interface{}, fields map[string]bool) map[string]interface{} {
	for column, secure := range fields {
		if !secure {
			delete(data, column)
		}
	}
	return data
}
