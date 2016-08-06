package dropbox

import (
	"fmt"
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
	Summary    string      `json:"error_summary"`
	Message    string      `json:"user_message"` // optionally present
	Err        interface{} `json:"error"`
}

// Error string.
func (e *Error) Error() string {
	if e.Summary != "" {
		return e.Summary
	}
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Status)
}

// Tag returns the inner tag for the error
func (e *Error) Tag() (tag, value string) {
	payload, ok := e.Err.(map[string]interface{})
	if !ok {
		val, ok := e.Err.(string)
		if ok {
			return "", val
		}

		return
	}

	tag, ok = payload[".tag"].(string)
	if !ok {
		return
	}

	data, ok := payload[tag].(map[string]interface{})
	if !ok {
		return
	}

	value, ok = data[".tag"].(string)
	if !ok {
		return
	}
	return
}
