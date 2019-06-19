package user

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"database/sql"

	mdlShared "github.com/agile-work/srv-mdl-shared"
	mdlSharedModels "github.com/agile-work/srv-mdl-shared/models"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"

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
	Token             string            `json:"token"`
	CreatedBy         string            `json:"created_by" sql:"created_by"`
	CreatedAt         time.Time         `json:"created_at" sql:"created_at"`
	UpdatedBy         string            `json:"updated_by" sql:"updated_by"`
	UpdatedAt         time.Time         `json:"updated_at" sql:"updated_at"`
}

// GetUserSelectableFields returns an array with the user fields
func GetUserSelectableFields() []string {
	return []string{
		constants.TableCoreUsers + ".id as user_id",
		constants.TableCoreUsers + ".username",
		constants.TableCoreUsers + ".first_name",
		constants.TableCoreUsers + ".last_name",
		constants.TableCoreUsers + ".email",
		constants.TableCoreUsers + ".language_code",
		constants.TableCoreUsers + ".receive_emails",
	}
}

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

// AuthCustomClaims used to parse token payload
type AuthCustomClaims struct {
	User User
	jwt.StandardClaims
}

// Users defines the array struct of this object
type Users []User

// Create persists the struct creating a new object in the database
func (u *User) Create(trs *db.Transaction, columns ...string) error {
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreUsers, u, columns...)
	if err != nil {
		return mdlShared.NewError("user create", err.Error())
	}
	u.ID = id

	// resource := localModels.Instance{}
	// resource.ID = sql.UUID()
	// resource.ParentID = user.ID
	// db.SetSchemaAudit(r, &resource)
	// trs.Add(sql.StructInsertQuery(shared.TableCustomResources, &resource, "", true))

	return nil
}

// LoadAll defines all instances from the object
func (u *Users) LoadAll(trs *db.Transaction, opt *db.Options) error {
	if err := db.SelectStructTx(trs.Tx, constants.TableCoreUsers, u, opt); err != nil {
		return mdlShared.NewError("users load", err.Error())
	}
	return nil
}

// Load defines only one object from the database
func (u *User) Load(trs *db.Transaction) error {
	if err := db.SelectStructTx(trs.Tx, constants.TableCoreUsers, u, &db.Options{
		Conditions: builder.Equal("username", u.Username),
	}); err != nil {
		return mdlShared.NewError("user load", err.Error())
	}
	return nil
}

// Update updates object data in the database
func (u *User) Update(trs *db.Transaction, columns []string, translations map[string]string) error {
	opt := &db.Options{Conditions: builder.Equal("username", u.Username)}

	if len(columns) > 0 {
		if err := db.UpdateStructTx(trs.Tx, constants.TableCoreUsers, u, opt, strings.Join(columns, ",")); err != nil {
			return mdlShared.NewError("user update", err.Error())
		}
	}

	if len(translations) > 0 {
		statement := builder.Update(constants.TableCoreUsers)
		for col, val := range translations {
			statement.JSON(col, mdlSharedModels.TranslationFieldsRequestLanguageCode)
			jsonVal, _ := json.Marshal(val)
			statement.Values(jsonVal)
		}
		statement.Where(opt.Conditions)
		if _, err := trs.Query(statement); err != nil {
			return mdlShared.NewError("user update", err.Error())
		}
	}

	return nil
}

// Delete deletes object from the database
func (u *User) Delete(trs *db.Transaction) error {
	if err := db.DeleteStructTx(trs.Tx, constants.TableCoreUsers, &db.Options{
		Conditions: builder.Equal("username", u.Username),
	}); err != nil {
		return mdlShared.NewError("user delete", err.Error())
	}
	return nil
}

// // LoadAllUsersByGroup return all instances from the object
// func LoadAllUsersByGroup(r *http.Request) *mdlShared.Response {
// 	viewGroupUsers := []models.ViewGroupUser{}
// 	groupID := chi.URLParam(r, "group_id")
// 	groupIDColumn := fmt.Sprintf("%s.group_id", shared.ViewCoreGroupUsers)
// 	condition := builder.Equal(groupIDColumn, groupID)

// 	return db.Load(r, &viewGroupUsers, "LoadAllUsersByGroup", shared.ViewCoreGroupUsers, condition)
// }

// // LoadAllPermissionsByGroup return all instances from the object
// func LoadAllPermissionsByGroup(r *http.Request) *mdlShared.Response {
// 	viewUserAllPermissions := []models.ViewUserAllPermissions{}
// 	userID := chi.URLParam(r, "user_id")
// 	userIDColumn := fmt.Sprintf("%s.user_id", shared.ViewCoreUserAllPermissions)
// 	languageCode := r.Header.Get("Content-Language")
// 	languageCodeColumn := fmt.Sprintf("%s.language_code", shared.ViewCoreUserAllPermissions)
// 	condition := builder.And(
// 		builder.Equal(userIDColumn, userID),
// 		builder.Equal(languageCodeColumn, languageCode),
// 	)

// 	return db.Load(r, &viewUserAllPermissions, "LoadAllPermissionsByGroup", shared.ViewCoreUserAllPermissions, condition)
// }

// Login validate credentials and return user token
func (u *User) Login(trs *db.Transaction) error {
	// response := &mdlShared.Response{
	// 	Code: http.StatusOK,
	// }
	// body, _ := ioutil.ReadAll(r.Body)

	// jsonMap := make(map[string]interface{})
	// err := json.Unmarshal(body, &jsonMap)
	// if err != nil {
	// 	response.Code = http.StatusInternalServerError
	// 	response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorParsingRequest, "Login unmarshal body", err.Error()))
	// 	return response
	// }

	if u.Email == "" || u.Password == "" {
		return mdlShared.NewError("user login", "invalid credentials body")
	}

	password := u.Password
	if err := db.SelectStructTx(trs.Tx, constants.TableCoreUsers, u, &db.Options{
		Conditions: builder.Equal("email", u.Email),
	}); err != nil {
		return mdlShared.NewError("user login load user", err.Error())
	}

	if u.ID == "" {
		return mdlShared.NewError("user login", "invalid email")
	}

	if u.Password != password {
		return mdlShared.NewError("user login", "invalid password")
	}

	u.Password = ""
	claims := AuthCustomClaims{
		*u,
		jwt.StandardClaims{
			Issuer: "API",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cryoSigningKey := []byte("AllYourBase") // TODO: Check the best place for this key, probably the config.toml
	tokenString, err := token.SignedString(cryoSigningKey)
	if err != nil {
		return mdlShared.NewError("user login generate token", err.Error())
	}

	u.Token = tokenString

	return nil
}

// // Load return user
// func (u *User) Load(userID string) error {
// 	userIDColumn := fmt.Sprintf("%s.id", constants.TableCoreUsers)
// 	activeColumn := fmt.Sprintf("%s.active", constants.TableCoreUsers)
// 	err := db.SelectStruct(constants.TableCoreUsers, u, &db.Options{Conditions: builder.And(
// 		builder.Equal(userIDColumn, userID),
// 		builder.Equal(activeColumn, true),
// 	)})
// 	if err != nil {
// 		return err
// 	}
// 	if !u.Active {
// 		return errors.New("Disabled user")
// 	}
// 	return nil
// }

// GetSecurityInstances return the initial statement to make a security query
func (u *User) GetSecurityInstances(schemaCode string) ([]map[string]interface{}, error) {
	statement, err := u.getSecurityQuery(schemaCode)
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

// GetSecurityQuery return the initial statement to make a security query
func (u *User) getSecurityQuery(schemaCode string) (*builder.Statement, error) {
	schemaTable := fmt.Sprintf("%s%s AS sch", constants.InstancesTablePrefix, schemaCode)
	securitySchema := u.Security.Schema[schemaCode]
	columnsJSON := strings.Split(securitySchema.QueryDataColumns, ",")
	columns := []string{}

	columns = append(columns, "sch.id")

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
				"sch.id = fld.schema_id",
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
