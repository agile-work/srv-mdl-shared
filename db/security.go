package db

import (
	"encoding/json"
	"fmt"

	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	sql "github.com/agile-work/srv-shared/sql-builder/db"
)

// StructurePermission defines attirbutes to a structure
type StructurePermission struct {
	ID                   string          `json:"id"`
	StructureCode        string          `json:"structure_code"`
	StructureType        string          `json:"structure_type"`  // field, widget, section...
	StructureClass       string          `json:"structure_class"` // field: lookup_tree, text, date, number, money | widget: chart, table
	StructureDefinitions json.RawMessage `json:"structure_definitions" field:"jsonb"`
	Permission           int             `json:"permission"`
	InstanceID           string          `json:"instance_id"`
}

// UserInstancePermission get user instace permissions
type UserInstancePermission struct {
	Scope                   string `json:"scope" sql:"scope"`
	UserID                  string `json:"user_id" sql:"user_id"`
	SchemaCode              string `json:"schema_code" sql:"schema_code"`
	TreeCode                string `json:"tree_code" sql:"tree_code"`
	TreeUnitCode            string `json:"tree_unit_code" sql:"tree_unit_code"`
	TreeUnitPath            string `json:"tree_unit_path" sql:"tree_unit_path"`
	TreeUnitPermissionScope string `json:"tree_unit_permission_scope" sql:"tree_unit_permission_scope"`
}

// GetUserAvailableFields returns a list of available structures based on user permission
// passing StructureType = "" returns all types of structures
func GetUserAvailableFields(userID, schemaCode, structureType, instanceID string) ([]StructurePermission, error) {
	statement := builder.Raw(`
		SELECT
			tab.user_id AS user_id,
			tab.schema_code AS schema_code,
			tab.id AS id,
			tab.structure_code AS structure_code,
			tab.structure_type AS structure_type,
			tab.structure_class AS structure_class,
			tab.structure_definitions AS structure_definitions,
			tab.instance_id AS instance_id,
			max(tab.permission) AS permission
		FROM (
			SELECT
				res.user_id AS user_id,
				sch.code AS schema_code,
				unit_path.structure_id AS id,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.code
					ELSE wdg.code
				END AS structure_code,
				unit_path.structure_type AS structure_type,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.field_type
					ELSE wdg.widget_type
				END AS structure_class,
				unit_path.permission AS permission,
				'unit' AS scope,
				fld.definitions AS structure_definitions,
				NULL AS instance_id
			FROM (
				SELECT
					res.parent_id AS user_id,
					trees->>'tree' AS tree,
					trees->>'tree_unit' AS tree_unit
				FROM
					cst_resources AS res,
					jsonb_array_elements(res.data->'trees') trees
				WHERE
					res.parent_id = '` + userID + `'
			) AS res
			JOIN
				core_trees tree
			ON
				tree.code = res.tree
			JOIN
				core_tree_units unit
			ON
				unit.tree_id = tree.id
				AND unit.code = res.tree_unit
			JOIN (
				SELECT 
					unit_path.tree_id AS tree_id,
					unit_path.permission_scope AS permission_scope,
					unit_path.path AS path,
					perm->>'structure_id' AS structure_id,
					perm->>'structure_type' AS structure_type,
					(perm->>'permission_type')::INT AS permission
				FROM 
					core_tree_units unit_path,
					jsonb_array_elements(unit_path.permissions) perm
				WHERE
					perm->>'structure_type' != 'schema'
			) AS unit_path
			ON
				unit_path.tree_id = unit.tree_id
			LEFT JOIN
				core_sch_fields fld
			ON
				unit_path.structure_id = fld.id
			LEFT JOIN
				core_schemas sch
			ON
				sch.id = fld.schema_id
			LEFT JOIN
				core_widgets wdg
			ON
				unit_path.structure_id = wdg.id
			WHERE
				sch.code = '` + schemaCode + `'
			AND
			(
				(
					unit_path.path = unit.path
					AND unit_path.permission_scope IS NOT NULL
				)
				OR
				(
					unit.path <@ unit_path.path
					AND unit_path.permission_scope = 'unit_and_descendent'
				)
			)

			UNION ALL

			SELECT
				res.parent_id AS user_id,
				sch.code AS schema_code,
				grp.structure_id AS id,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.code
					ELSE wdg.code
				END AS structure_code,
				grp.structure_type,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.field_type
					ELSE wdg.widget_type
				END AS structure_class,
				grp.permission,
				CASE
					WHEN grp.tree_unit_id IS NULL
					THEN 'group'
					ELSE 'unit_group'
				END AS scope,
				fld.definitions AS structure_definitions,
				NULL AS instance_id
			FROM cst_resources AS res
			JOIN (
				SELECT 
					grp.id AS id,
					grp.tree_unit_id AS tree_unit_id,
					grp.users AS users,
					perm->>'structure_id' AS structure_id,
					perm->>'structure_type' AS structure_type,
					(perm->>'permission_type')::INT AS permission
				FROM 
					core_groups grp,
					jsonb_array_elements(grp.permissions) perm
				WHERE
					perm->>'structure_type' != 'schema'
					AND grp.users @> ('[{"id":"` + userID + `"}]')::JSONB
			) AS grp
			ON
				res.parent_id = '` + userID + `'
				AND grp.users @> ('[{"id":"' || res.parent_id || '"}]')::JSONB
			LEFT JOIN
				core_sch_fields fld
			ON
				grp.structure_id = fld.id
			LEFT JOIN
				core_schemas sch
			ON
				sch.id = fld.schema_id
			LEFT JOIN
				core_widgets wdg
			ON
				grp.structure_id = wdg.id
			WHERE
				sch.code = '` + schemaCode + `'

			UNION ALL

			SELECT
				inst.user_id AS user_id,
				inst.instance_type AS schema_code,
				inst.structure_id AS id,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.code
					ELSE wdg.code
				END AS structure_code,
				inst.structure_type AS structure_type,
				CASE
					WHEN fld.id IS NOT NULL THEN fld.field_type
					ELSE wdg.widget_type
				END AS structure_class,
				inst.permission AS permission,
				'instance_' || inst.source_type AS scope,
				fld.definitions AS structure_definitions,
				inst.instance_id AS instance_id
			FROM (
				SELECT 
					inst.user_id AS user_id,
					inst.instance_type AS instance_type,
					inst.source_type AS source_type,
					perm->>'structure_id' AS structure_id,
					perm->>'structure_type' AS structure_type,
					(perm->>'permission_type')::INT AS permission,
					inst.instance_id AS instance_id
				FROM 
					core_instance_premissions inst,
					jsonb_array_elements(inst.permissions) perm
				WHERE
					perm->>'structure_type' != 'schema'
					AND inst.instance_type = '` + schemaCode + `'
					AND (
						'` + instanceID + `' = '' 
						OR inst.instance_id = '` + instanceID + `'
					)
			) AS inst    
			LEFT JOIN
				core_sch_fields fld
			ON
				inst.structure_id = fld.id
			LEFT JOIN
				core_widgets wdg
			ON
				inst.structure_id = wdg.id
		) AS tab
		GROUP BY
			tab.user_id,
			tab.schema_code,
			tab.id,
			tab.structure_code,
			tab.structure_type,
			tab.structure_class,
			tab.structure_definitions,
			tab.instance_id
	`)

	rows, err := sql.Query(statement)
	if err != nil {
		return nil, err
	}

	permission := []StructurePermission{}
	err = sql.StructScan(rows, &permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// GetTreeSecurityFieldsFromSchema get tree security fields from schema to make a condition to get instances
func GetTreeSecurityFieldsFromSchema(schemaCode string) ([]map[string]interface{}, error) {
	statement := builder.Select(
		"fld.code",
		"fld.definitions->'definitions'->>'tree' AS tree",
	).From(
		fmt.Sprintf("%s AS sch", shared.TableCoreSchemas),
	).Join(
		fmt.Sprintf("%s AS fld", shared.TableCoreSchemaFields),
		"fld.schema_id = sch.id",
	).Where(
		builder.And(
			builder.Equal("sch.code", schemaCode),
			builder.Equal("fld.field_type", "lookup"),
			builder.Raw(`fld.definitions @> '{"definitions": {"lookup_type": "tree", "is_security": true}}'`),
		),
	)

	rows, err := sql.Query(statement)
	if err != nil {
		return nil, err
	}

	treeFields, err := sql.MapScan(rows)
	if err != nil {
		return nil, err
	}

	return treeFields, nil
}

// GetUserInstancePermissions get all instance permissions for a schema and user,
// except from the instance permission table.
func GetUserInstancePermissions(userID, schemaCode string) ([]UserInstancePermission, error) {
	statement := builder.Raw(`
		SELECT
			'group' AS scope,
			res.parent_id AS user_id,
			sch.code AS schema_code,
			NULL AS tree_code,
			NULL AS tree_unit_code,
			NULL::TEXT AS tree_unit_path,
			NULL AS tree_unit_permission_scope
		FROM
		(
			SELECT
				grp.id AS id,
				grp.tree_unit_id AS tree_unit_id,
				grp.users AS users,
				perm->>'structure_id' AS structure_id,
				perm->>'structure_type' AS structure_type,
				(perm->>'permission_type')::INT AS permission_type
			FROM
				core_groups AS grp,
				jsonb_array_elements(grp.permissions) AS perm
			WHERE
				(perm->>'permission_type')::INT = 100
				AND grp.tree_unit_id IS NULL
				AND grp.active = TRUE
				AND grp.users @> ('[{"id":"` + userID + `"}]')::JSONB
		) AS grp
		JOIN
			cst_resources AS res
		ON
			res.parent_id = '` + userID + `'
			AND grp.users @> ('[{"id":"' || res.parent_id || '"}]')::JSONB
		JOIN
			core_schemas AS sch
		ON
			sch.id = grp.structure_id
			AND sch.code = '` + schemaCode + `'

		UNION ALL

		SELECT
			'group_unit' AS scope,
			res.user_id AS user_id,
			sch.code AS schema_code,
			res.tree_code AS tree_code,
			res.tree_unit_code AS tree_unit_code,
			unit_res.path::TEXT AS tree_unit_path,
			min(grp.tree_unit_permission_scope) AS tree_unit_permission_scope
		FROM
		(
			SELECT
				grp.id AS id,
				grp.tree_unit_id AS tree_unit_id,
				grp.tree_unit_permission_scope AS tree_unit_permission_scope,
				grp.users AS users,
				perm->>'structure_id' AS structure_id,
				perm->>'structure_type' AS structure_type,
				(perm->>'permission_type')::INT AS permission_type
			FROM
				core_groups AS grp,
				jsonb_array_elements(grp.permissions) AS perm
			WHERE
				(perm->>'permission_type')::INT = 100
				AND grp.tree_unit_id IS NOT NULL
				AND grp.active = TRUE				
				AND grp.users @> ('[{"id":"` + userID + `"}]')::JSONB
		) AS grp
		JOIN
			core_tree_units AS unit
		ON
			unit.id = grp.tree_unit_id
			AND unit.active = TRUE
		JOIN
			core_trees AS trees
		ON
			trees.id = unit.tree_id
		JOIN
			(
				SELECT
					res.parent_id AS user_id,
					trees->>'tree' AS tree_code,
					trees->>'tree_unit' AS tree_unit_code
				FROM
					cst_resources AS res,
					jsonb_array_elements(res.data->'trees') trees
				WHERE
					res.parent_id = '` + userID + `'
			) AS res
		ON
			grp.users @> ('[{"id":"' || res.user_id || '"}]')::JSONB
			AND res.tree_code = trees.code
		JOIN
			core_trees AS tree_res
		ON
			tree_res.code = res.tree_code
		JOIN
			core_tree_units AS unit_res
		ON
			unit_res.tree_id = tree_res.id
			AND unit_res.code = res.tree_unit_code
		JOIN
			core_schemas AS sch
		ON
			sch.id = grp.structure_id
			AND sch.code = '` + schemaCode + `'
		GROUP BY
			res.user_id,
			sch.code,
			res.tree_code,
			res.tree_unit_code,
			unit_res.path

		UNION ALL

		SELECT
			'unit' AS scope,
			res.user_id AS user_id,
			sch.code AS schema_code,
			res.tree_code AS tree_code,
			res.tree_unit_code AS tree_unit_code,
			unit.path::TEXT AS tree_unit_path,
			min(unit_path.permission_scope) AS tree_unit_permission_scope
		FROM
			(
				SELECT
					res.parent_id AS user_id,
					trees->>'tree' AS tree_code,
					trees->>'tree_unit' AS tree_unit_code
				FROM
					cst_resources AS res,
					jsonb_array_elements(res.data->'trees') trees
				WHERE
					res.parent_id = '` + userID + `'
			) AS res
		JOIN
			core_trees AS tree
		ON
			tree.code = res.tree_code
		JOIN
			core_tree_units AS unit
		ON
			unit.tree_id = tree.id
			AND unit.code = res.tree_unit_code
		JOIN (
			SELECT
				unit_path.id AS unit_id,
				unit_path.tree_id AS tree_id,
				unit_path.permission_scope AS permission_scope,
				unit_path.path AS path,
				perm->>'structure_id' AS structure_id,
				perm->>'structure_type' AS structure_type,
				(perm->>'permission_type')::INT AS permission_type
			FROM 
				core_tree_units AS unit_path,
				jsonb_array_elements(unit_path.permissions) AS perm
			WHERE
				(perm->>'permission_type')::INT = 100
				AND unit_path.active = TRUE
		) AS unit_path
		ON
			unit_path.tree_id = unit.tree_id
		JOIN
			core_schemas AS sch
		ON
			sch.id = unit_path.structure_id
			AND sch.code = '` + schemaCode + `'
		WHERE
			(
				unit_path.unit_id = unit.id
				AND unit_path.permission_scope IS NOT NULL
			)
			OR
			(
				unit.path <@ unit_path.path
				AND unit_path.permission_scope = 'unit_and_descendent'
			)
		GROUP BY
			res.user_id,
			sch.code,
			res.tree_code,
			res.tree_unit_code,
			unit.path
	`)

	rows, err := sql.Query(statement)
	if err != nil {
		return nil, err
	}

	userInstancePermissions := []UserInstancePermission{}
	err = sql.StructScan(rows, &userInstancePermissions)
	if err != nil {
		return nil, err
	}

	return userInstancePermissions, nil
}
