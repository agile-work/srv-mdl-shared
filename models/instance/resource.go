package instance

// // LoadAllResources return all instances from the object
// func LoadAllResources(r *http.Request) *mdlShared.Response {
// 	response := &mdlShared.Response{
// 		Code: http.StatusOK,
// 	}
// 	userID := r.Header.Get("userID")
// 	securityFields, err := db.GetUserAvailableFields(userID, "resources", shared.SecurityStructureField, "")
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "Loading user fields permission", err.Error()))
// 		return response
// 	}
// 	fields := []string{}
// 	treeJoin := make(map[string]string)
// 	columns := []string{}

// 	for _, f := range securityFields {
// 		fields = append(fields, f.StructureCode)
// 	}

// 	columns = append(columns, models.GetUserSelectableFields()...)
// 	columns = append(columns, shared.TableCustomResources+".id")

// 	on := fmt.Sprintf("%s.id = %s.parent_id", shared.TableCoreUsers, shared.TableCustomResources)
// 	statement := builder.Select(columns...).JSON("data", fields...).From(shared.TableCustomResources).Join(shared.TableCoreUsers, on)
// 	for table, on := range treeJoin {
// 		statement.Join(table, on)
// 	}

// 	rows, err := sql.Query(statement)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadAllResources", err.Error()))
// 		return response
// 	}

// 	results, err := sql.MapScan(rows)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadAllResources Parsing query rows to map", err.Error()))
// 		return response
// 	}

// 	resp.Data = results
// 	return response
// }

// // LoadResource return one instance from the object
// func LoadResource(r *http.Request) *mdlShared.Response {
// 	response := &mdlShared.Response{
// 		Code: http.StatusOK,
// 	}
// 	userID := r.Header.Get("userID")
// 	resourceID := chi.URLParam(r, "resource_id")

// 	securityFields, err := db.GetUserAvailableFields(userID, "resources", shared.SecurityStructureField, resourceID)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "Loading user fields permission", err.Error()))
// 		return response
// 	}
// 	fields := []string{}
// 	treeJoin := make(map[string]string)
// 	columns := []string{}

// 	for _, f := range securityFields {
// 		fields = append(fields, f.StructureCode)
// 	}

// 	columns = append(columns, models.GetUserSelectableFields()...)
// 	columns = append(columns, shared.TableCustomResources+".id")

// 	on := fmt.Sprintf("%s.id = %s.parent_id", shared.TableCoreUsers, shared.TableCustomResources)
// 	statement := builder.Select(columns...).JSON("data", fields...).From(shared.TableCustomResources).Join(shared.TableCoreUsers, on)
// 	for table, on := range treeJoin {
// 		statement.Join(table, on)
// 	}
// 	statement.Where(builder.Equal("id", resourceID))

// 	rows, err := sql.Query(statement)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadResource", err.Error()))
// 		return response
// 	}

// 	results, err := sql.MapScan(rows)
// 	if err != nil {
// 		response.Code = http.StatusInternalServerError
// 		response.Errors = append(response.Errors, mdlShared.NewResponseError(shared.ErrorLoadingInstances, "LoadResource Parsing query rows to map", err.Error()))
// 		return response
// 	}

// 	if len(results) > 0 {
// 		resp.Data = results[0]
// 	}
// 	return response
// }

// // UpdateResource update an instance from the object
// func UpdateResource(r *http.Request) *mdlShared.Response {
// 	resourceID := chi.URLParam(r, "resource_id")

// 	resourceMap := map[string]interface{}{}

// 	response := db.GetResponse(r, &resourceMap, "UpdateResource")
// 	if response.Code != http.StatusOK {
// 		return response
// 	}

// 	// TODO: validate resource fields at resourceMap before update

// 	sql.UpdateStructToJSON(resourceID, "data", "cst_resources", &resourceMap)
// 	return response
// }
