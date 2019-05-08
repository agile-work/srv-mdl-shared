package shared

//Metadata defines metadata on the response
type Metadata struct {
}

// ResponseError defines the struct to the api response error
type ResponseError struct {
	Code  string `json:"code"`
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

// NewResponseError defines a structure to encode api response data
func NewResponseError(code string, scope, err string) ResponseError {
	return ResponseError{
		Code:  code,
		Scope: scope,
		Error: err,
	}
}
