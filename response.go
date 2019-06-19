package shared

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"

	"github.com/agile-work/srv-shared/util"

	"github.com/go-chi/render"
)

// Metadata defines metadata on the response
type Metadata struct {
	Filter     map[string]map[string]interface{} `json:"filter,omitempty"`
	Pagination map[string]int                    `json:"pagination,omitempty"`
	Order      []map[string]string               `json:"order,omitempty"`
	Columns    string                            `json:"columns,omitempty"`
	Group      string                            `json:"group,omitempty"`
}

// Load gets metadata from request
func (m *Metadata) Load(req *http.Request) error {
	query := req.URL.Query()
	metaDataStr := query.Get("metadata")
	metaDataBytes, err := json.Marshal(metaDataStr)
	if err != nil {
		return NewError("metadata load marshal", err.Error())
	}
	if err := json.Unmarshal(metaDataBytes, m); err != nil {
		return NewError("metadata load unmarshal", err.Error())
	}
	return nil
}

// GenerateDBOptions convert the metadata in a db options
func (m *Metadata) GenerateDBOptions() *db.Options {
	opt := &db.Options{}
	opt.Columns = strings.Split(m.Columns, ",")
	opt.Limit = m.Pagination["totalItens"]
	opt.Offset = (m.Pagination["selected"] * m.Pagination["totalItens"]) - m.Pagination["totalItens"]
	for _, row := range m.Order {
		for column, order := range row {
			if order == "asc" {
				opt.AddOrderBy(builder.Asc(column))
			} else {
				opt.AddOrderBy(builder.Desc(column))
			}
		}
	}
	for column, prop := range m.Filter {
		for operator, value := range prop {
			switch operator {
			case "=":
				opt.AddCondition(builder.Equal(column, value))
				break
			case "!=":
				opt.AddCondition(builder.NotEqual(column, value))
				break
			case ">":
				opt.AddCondition(builder.GreaterThen(column, value))
				break
			case ">=":
				opt.AddCondition(builder.GreaterOrEqual(column, value))
				break
			case "<":
				opt.AddCondition(builder.LowerThen(column, value))
				break
			case "<=":
				opt.AddCondition(builder.LowerOrEqual(column, value))
				break
			}
		}
	}
	return opt
}

// ResponseError defines the struct to the api response error
type ResponseError struct {
	Scope string `json:"scope"`
	Error string `json:"erro"`
}

// Response defines the struct to the api response
type Response struct {
	Code     int             `json:"code"`
	Metadata Metadata        `json:"metadata"`
	Data     interface{}     `json:"data"`
	Errors   []ResponseError `json:"errors"`
}

// NewError creats a new error in response
func (r *Response) NewError(code int, scope, err string) {
	r.Code = code
	r.Errors = append(r.Errors, ResponseError{
		Scope: scope,
		Error: err,
	})
}

// Render return a http response
func (r *Response) Render(res http.ResponseWriter, req *http.Request) {
	render.Status(req, r.Code)
	render.JSON(res, req, r)
}

// Load get request body to object and creates a response
func (r *Response) Load(req *http.Request, object interface{}) error {
	r.Code = http.StatusOK
	body, _ := util.GetBody(req)
	if len(body) > 0 {
		err := json.Unmarshal(body, object)
		if err != nil {
			return NewError("response load unmarshal body", err.Error())
		}
		if req.Method == http.MethodPost {
			err = Validate.Struct(object)
			if err != nil {
				return NewError("response load invalid body", err.Error())
			}
		}
	}

	util.SetSchemaAudit(req, object)
	return nil
}

// NewResponse make a new response
func NewResponse() *Response {
	return &Response{Code: http.StatusOK}
}
