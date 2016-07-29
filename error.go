package dropbox

import (
	"net/http"
)

// error tag constant values
const (
	TooManyRequests = "too_many_requests"
)

// Error response.
type Error struct {
	Status     string
	StatusCode int
	Header     http.Header
	Summary    string                 `json:"error_summary"`
	Message    string                 `json:"user_message"` // optionally present
	Err        map[string]interface{} `json:"error"`
}

// Error string.
func (e *Error) Error() string {
	return e.Summary
}

// Tag returns the inner tag for the error
func (e *Error) Tag() (tag, value string) {
	var ok bool
	tag, ok = e.Err[".tag"].(string)
	if !ok {
		return
	}

	data, ok := e.Err[tag].(map[string]interface{})
	if !ok {
		return
	}

	value, ok = data[".tag"].(string)
	if !ok {
		return
	}
	return
}
