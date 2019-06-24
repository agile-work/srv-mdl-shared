package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	shared "github.com/agile-work/srv-mdl-shared"
	"github.com/agile-work/srv-mdl-shared/models/customerror"

	"github.com/agile-work/srv-mdl-shared/util"

	"github.com/go-chi/render"
)

// Response defines the struct to the api response
type Response struct {
	Code     int         `json:"code"`
	Metadata Metadata    `json:"metadata,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    error       `json:"error,omitempty"`
}

// Render return a http response
func (r *Response) Render(res http.ResponseWriter, req *http.Request) {
	render.Status(req, r.Code)
	render.JSON(res, req, r)
}

// NewError creats a new error in response
func (r *Response) NewError(scope string, err error) {
	if custom, ok := err.(*customerror.Error); ok {
		custom.Scope = fmt.Sprintf("%s - %s", scope, custom.Scope)
		r.Code = custom.Code
		r.Error = err
	} else {
		r.Code = http.StatusInternalServerError
		msg := fmt.Sprintf("%s - %s", scope, err.Error())
		r.Error = errors.New(msg)
	}
}

// Parse get request body to object and creates a response
func (r *Response) Parse(req *http.Request, object interface{}) error {
	r.Code = http.StatusOK
	body, _ := util.GetBody(req)
	if len(body) > 0 {
		err := json.Unmarshal(body, object)
		if err != nil {
			return customerror.New(http.StatusBadRequest, "response load unmarshal body", err.Error())
		}
		if req.Method == http.MethodPost {
			err = shared.Validate.Struct(object)
			if err != nil {
				return customerror.New(http.StatusBadRequest, "response load invalid body", err.Error())
			}
		}
	}

	util.SetSchemaAudit(req, object)
	return nil
}

// New make a new response
func New() *Response {
	return &Response{Code: http.StatusOK}
}