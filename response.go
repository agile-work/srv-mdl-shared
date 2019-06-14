package shared

import (
	"net/http"

	"github.com/go-chi/render"
)

// Metadata defines metadata on the response
type Metadata struct {
}

// ResponseError defines the struct to the api response error
type ResponseError struct {
	Code  string `json:"code,omitempty"` // TODO: Deprecated
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

// Render
func (r *Response) Render(res http.ResponseWriter, req *http.Request) {
	render.Status(req, r.Code)
	render.JSON(res, req, r)
}

// TODO: Deprecated
// NewResponseError defines a structure to encode api response data
func NewResponseError(code string, scope, err string) ResponseError {
	return ResponseError{
		Code:  code,
		Scope: scope,
		Error: err,
	}
}
