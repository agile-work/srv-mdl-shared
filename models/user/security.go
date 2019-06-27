package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agile-work/srv-shared/util"

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
func (u *User) GetSecurityInstances(schemaCode string, opt *db.Options) ([]map[string]interface{}, error) {
	// TODO: Implementar o db options nas queries
	statement, err := u.GetSecurityQuery(schemaCode)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}

	results, err := u.securityMapScan(schemaCode, rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetSecurityQueryWithSub return the initial statement to make a security query
func (u *User) GetSecurityQueryWithSub(schemaCode string, subSelect interface{}, opt *db.Options) (*builder.Statement, error) {
	schemaTable := fmt.Sprintf("%s%s AS sch", constants.InstancesTablePrefix, schemaCode)
	securitySchema := u.Security.Schema[schemaCode]
	columns := []string{}

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree := range securitySchema.SecurityFields {
			columns = append(columns, fmt.Sprintf("unit_%s.path::TEXT AS unit_%s", schemaColumnTree, schemaColumnTree))
		}
	}

	columns = append(columns, "sub.*")

	statement := builder.Select(columns...).From(schemaTable).Join(
		fmt.Sprintf("(%s) sub", subSelect),
		"sub.id = sch.id",
	)

	u.loadSecurityTreeJoins(schemaCode, statement)
	u.loadSecurityConditions(schemaCode, statement)

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

// GetSecurityQuery return the initial statement to make a security query
func (u *User) GetSecurityQuery(schemaCode string) (*builder.Statement, error) {
	schemaTable := fmt.Sprintf("%s%s AS sch", constants.InstancesTablePrefix, schemaCode)
	securitySchema := u.Security.Schema[schemaCode]
	columnsJSON := strings.Split(securitySchema.QueryDataColumns, ",")
	columns := []string{
		"sch.id",
	}

	if securitySchema.PermissionInstance == "custom" {
		for schemaColumnTree := range securitySchema.SecurityFields {
			columns = append(columns, fmt.Sprintf("unit_%s.path::TEXT AS unit_%s", schemaColumnTree, schemaColumnTree))
		}
	}

	if columnsJSON[0] == "*" {
		columnsJSON = []string{}

		// TODO: Pensar em um jeito de trocar essa consulta de fields por um cache com REDIS
		rows, err := db.Query(
			builder.Select(
				"fld.code",
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
	}

	statement := builder.Select(columns...).JSON("data", columnsJSON...).From(schemaTable)

	u.loadSecurityTreeJoins(schemaCode, statement)
	u.loadSecurityConditions(schemaCode, statement)

	return statement, nil
}

// loadSecurityTreeJoins put the tree joins on statement to make a security query
func (u *User) loadSecurityTreeJoins(schemaCode string, statementSchema *builder.Statement) {
	securitySchema := u.Security.Schema[schemaCode]

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
func (u *User) loadSecurityConditions(schemaCode string, statementSchema *builder.Statement) {
	securitySchema := u.Security.Schema[schemaCode]
	securityInstanceSchema := u.SecurityInstances.Schema[schemaCode]
	conditions := []builder.Builder{}

	if securitySchema.PermissionInstance == "custom" {
		// TODO: Usar o SecurityFields em um cache com REDIS no lugar da informação no security do usuário
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

// securityMapScan checks the security in the instance columns and clears if user do not have permission
func (u *User) securityMapScan(schemaCode string, rows *sql.Rows) ([]map[string]interface{}, error) {
	securitySchema := u.Security.Schema[schemaCode]
	securityInstanceSchema := u.SecurityInstances.Schema[schemaCode]
	deleteColumn := true
	requiredFields := map[string]bool{}
	requiredFields["id"] = true
	// TODO: Usar o SecurityFields em um cache com REDIS no lugar da informação no security do usuário
	for schemaColumnTree := range securitySchema.SecurityFields {
		// TODO: Usar o requiredFields em um cache com REDIS
		requiredFields[schemaColumnTree] = true
		requiredFields["unit_"+schemaColumnTree] = true
	}
	cols, _ := rows.Columns()
	results := []map[string]interface{}{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		mapJSON := make(map[string]interface{})
		for i, column := range cols {
			val := columnPointers[i].(*interface{})
			if source, ok := columns[i].([]byte); ok {
				var raw json.RawMessage
				_ = json.Unmarshal(source, &raw)
				mapJSON[column] = raw
			} else {
				mapJSON[column] = *val
			}
			if securitySchema.PermissionStructure == "custom" {
				if requiredFields[column] {
					continue
				}
				deleteColumn = true
				if instanceOptions, ok := securityInstanceSchema.Instance[mapJSON["id"].(string)]; !ok || (ok && instanceOptions.PermissionScope != "replace") {
					if securitySchema.PermissionStructure == "custom" {
						if rule, ok := securitySchema.Edit[column]; ok {
							deleteColumn = !rule
						}
						if rule, ok := securitySchema.View[column]; ok {
							deleteColumn = !rule
						}
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
									instanceSecurityUnitPath := mapJSON[("unit_" + schemaColumnTree)].(string)
									if (securityUnitPath[len(securityUnitPath)-1:] == "*" &&
										(instanceSecurityUnitPath == userUnit ||
											strings.HasPrefix(instanceSecurityUnitPath, userUnit+".") ||
											strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit) ||
											strings.Contains(instanceSecurityUnitPath, "."+userUnit+"."))) ||
										(securityUnitPath[len(securityUnitPath)-1:] != "*" &&
											(instanceSecurityUnitPath == userUnit ||
												strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit))) {
										if securitySchemaTree.Permission != "custom" {
											deleteColumn = false
										}
										if rule, ok := securitySchemaTree.Edit[column]; ok {
											deleteColumn = !rule
										}
										if rule, ok := securitySchemaTree.View[column]; ok {
											deleteColumn = !rule
										}
										break LoopTrees
									}
								}
							}
						}
					} else {
						deleteColumn = false
					}
				}

				if instanceOptions, ok := securityInstanceSchema.Instance[mapJSON["id"].(string)]; ok {
					if instanceOptions.PermissionScope == "replace" {
						deleteColumn = true
					}
					if instanceOptions.Permission != "custom" {
						deleteColumn = false
					}
					if rule, ok := instanceOptions.Edit[column]; ok {
						deleteColumn = !rule
					}
					if rule, ok := instanceOptions.View[column]; ok {
						deleteColumn = !rule
					}
				}
				if deleteColumn {
					delete(mapJSON, column)
				}
			}
		}
		for schemaColumnTree := range securitySchema.SecurityFields {
			delete(mapJSON, "unit_"+schemaColumnTree)
		}
		results = append(results, mapJSON)
	}
	rows.Close()

	return results, nil
}

// SecurityMapScanWithFields checks the security in the instance columns and clears if user do not have permission
func (u *User) SecurityMapScanWithFields(schemaCode string, rows *sql.Rows, opt *db.Options, fields map[string]map[string]string) ([]map[string]interface{}, error) {
	securitySchema := u.Security.Schema[schemaCode]
	securityInstanceSchema := u.SecurityInstances.Schema[schemaCode]
	deleteColumn := true
	requiredFields := map[string]bool{}
	requiredFields["id"] = true
	// TODO: Usar o SecurityFields em um cache com REDIS no lugar da informação no security do usuário
	for schemaColumnTree := range securitySchema.SecurityFields {
		// TODO: Usar o requiredFields em um cache com REDIS
		requiredFields[schemaColumnTree] = true
		requiredFields["unit_"+schemaColumnTree] = true
	}
	cols, _ := rows.Columns()
	results := []map[string]interface{}{}
	schemaFields := fields[schemaCode]

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		mapJSON := make(map[string]interface{})
		for i, column := range cols {
			val := columnPointers[i].(*interface{})
			if source, ok := columns[i].([]byte); ok {
				var raw json.RawMessage
				_ = json.Unmarshal(source, &raw)
				mapJSON[column] = raw
			} else {
				mapJSON[column] = *val
			}
			if requiredFields[column] {
				continue
			}
			if !util.Contains(opt.Columns, column) {
				delete(mapJSON, column)
				continue
			}
			if securitySchema.PermissionStructure == "custom" {
				if securityColumn, ok := schemaFields[column]; ok {
					deleteColumn = true
					if instanceOptions, ok := securityInstanceSchema.Instance[mapJSON["id"].(string)]; !ok || (ok && instanceOptions.PermissionScope != "replace") {
						if securitySchema.PermissionStructure == "custom" {
							if rule, ok := securitySchema.Edit[securityColumn]; ok {
								deleteColumn = !rule
							}
							if rule, ok := securitySchema.View[securityColumn]; ok {
								deleteColumn = !rule
							}
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
										instanceSecurityUnitPath := mapJSON[("unit_" + schemaColumnTree)].(string)
										if (securityUnitPath[len(securityUnitPath)-1:] == "*" &&
											(instanceSecurityUnitPath == userUnit ||
												strings.HasPrefix(instanceSecurityUnitPath, userUnit+".") ||
												strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit) ||
												strings.Contains(instanceSecurityUnitPath, "."+userUnit+"."))) ||
											(securityUnitPath[len(securityUnitPath)-1:] != "*" &&
												(instanceSecurityUnitPath == userUnit ||
													strings.HasSuffix(instanceSecurityUnitPath, "."+userUnit))) {
											if securitySchemaTree.Permission != "custom" {
												deleteColumn = false
											}
											if rule, ok := securitySchemaTree.Edit[securityColumn]; ok {
												deleteColumn = !rule
											}
											if rule, ok := securitySchemaTree.View[securityColumn]; ok {
												deleteColumn = !rule
											}
											break LoopTrees
										}
									}
								}
							}
						} else {
							deleteColumn = false
						}
					}

					if instanceOptions, ok := securityInstanceSchema.Instance[mapJSON["id"].(string)]; ok {
						if instanceOptions.PermissionScope == "replace" {
							deleteColumn = true
						}
						if instanceOptions.Permission != "custom" {
							deleteColumn = false
						}
						if rule, ok := instanceOptions.Edit[securityColumn]; ok {
							deleteColumn = !rule
						}
						if rule, ok := instanceOptions.View[securityColumn]; ok {
							deleteColumn = !rule
						}
					}
					if deleteColumn {
						delete(mapJSON, column)
					}
				} else {
					continue
				}
			}
		}
		for schemaColumnTree := range securitySchema.SecurityFields {
			delete(mapJSON, "unit_"+schemaColumnTree)
		}
		results = append(results, mapJSON)
	}
	rows.Close()

	return results, nil
}
