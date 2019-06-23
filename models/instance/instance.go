package instance

import (
	"encoding/json"
	"time"
)

// Instance defines the struct of this object
type Instance struct {
	ID        string          `json:"id" sql:"id" pk:"true"`
	ParentID  string          `json:"parent_id" sql:"parent_id" pk:"true"`
	Data      json.RawMessage `json:"data" sql:"data" field:"jsonb"`
	CreatedBy string          `json:"created_by" sql:"created_by"`
	CreatedAt time.Time       `json:"created_at" sql:"created_at"`
	UpdatedBy string          `json:"updated_by" sql:"updated_by"`
	UpdatedAt time.Time       `json:"updated_at" sql:"updated_at"`
}

// EntityInstancePermission defines the struct of this object
type EntityInstancePermission struct {
	ID           string               `json:"id" sql:"id" pk:"true"`
	UserID       string               `json:"user_id" sql:"user_id"`
	SourceType   string               `json:"source_type" sql:"source_type"`
	SourceID     string               `json:"source_id" sql:"source_id"`
	InstanceID   string               `json:"instance_id" sql:"instance_id"`
	InstanceType string               `json:"instance_type" sql:"instance_type"`
	Permissions  []InstancePermission `json:"permissions" sql:"permissions" field:"jsonb"`
	CreatedBy    string               `json:"created_by" sql:"created_by"`
	CreatedAt    time.Time            `json:"created_at" sql:"created_at"`
	UpdatedBy    string               `json:"updated_by" sql:"updated_by"`
	UpdatedAt    time.Time            `json:"updated_at" sql:"updated_at"`
}

// InstancePermission defines the struct of this object
type InstancePermission struct {
	ID             string `json:"id" sql:"id" pk:"true"`
	StructureType  string `json:"structure_type" sql:"structure_type"`
	StructureID    string `json:"structure_id" sql:"structure_id"`
	PermissionType int    `json:"permission_type" sql:"permission_type"`
	ConditionQuery string `json:"condition_query" sql:"condition_query"`
}

// // CreateSchemaInstance persists the request body creating a new object in the database
// func CreateSchemaInstance(r *http.Request) *mdlShared.Response {
// 	instance := Instance{}
// 	schemaCode := chi.URLParam(r, "schema_code")

// 	return db.Create(r, &instance, "CreateInstance", fmt.Sprintf("%s%s", shared.InstancesTablePrefix, schemaCode))
// }

// // LoadAllSchemaInstances return all instances from the object
// func LoadAllSchemaInstances(r *http.Request) *mdlShared.Response {
// 	response := &mdlShared.Response{
// 		Code: http.StatusOK,
// 	}

// 	schemaCode := chi.URLParam(r, "schema_code")
// 	user := modulemdlSharedModels.User{}
// 	err := user.Load(r.Header.Get("userID"))
// 	if err != nil {
// 		response.Code = http.StatusForbidden
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadAllInstances Load user", err.Error()))
// 		return response
// 	}

// 	results, err := user.GetSecurityInstances(schemaCode)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadAllInstances LoadAllInstances loading instances", err.Error()))
// 		return response
// 	}

// 	resp.Data = results

// 	return response
// }

// // LoadSchemaInstance return only one object from the database
// func LoadSchemaInstance(r *http.Request) *mdlShared.Response {
// 	return nil
// }

// // UpdateSchemaInstance updates object data in the database
// func UpdateSchemaInstance(r *http.Request) *mdlShared.Response {
// 	schemaCode := chi.URLParam(r, "schema_code")
// 	instanceID := chi.URLParam(r, "instance_id")
// 	table := fmt.Sprintf("%s%s", shared.InstancesTablePrefix, schemaCode)
// 	instanceIDColumn := fmt.Sprintf("%s.id", table)
// 	condition := builder.Equal(instanceIDColumn, instanceID)

// 	instance := Instance{
// 		ID: instanceID,
// 	}

// 	return db.Update(r, &instance, "UpdateInstance", table, condition)
// }

// // DeleteSchemaInstance deletes object from the database
// func DeleteSchemaInstance(r *http.Request) *mdlShared.Response {
// 	schemaCode := chi.URLParam(r, "schema_code")
// 	instanceID := chi.URLParam(r, "instance_id")
// 	table := fmt.Sprintf("%s%s", shared.InstancesTablePrefix, schemaCode)
// 	instanceIDColumn := fmt.Sprintf("%s.id", table)
// 	condition := builder.Equal(instanceIDColumn, instanceID)

// 	return db.Remove(r, "DeleteInstance", table, condition)
// }
