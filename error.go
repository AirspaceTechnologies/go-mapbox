package mapbox

import "fmt"

type MapboxError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error"`
}

////////////////////////////////////////////////////////////////////////////////

func NewMapboxError(statusCode int, message string) MapboxError {
	return MapboxError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e MapboxError) Error() string {
	return fmt.Sprintf("api error(%v): %v", e.StatusCode, e.Message)
}
