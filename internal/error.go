package internal

import (
	"encoding/json"

	saved "github.com/savedhq/sdk-go"
)

// APIError is a structured error from the API.
type APIError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return e.Message
}

// PrintAPIError formats an SDK error.
func PrintAPIError(err error) error {
	if apiErr, ok := err.(*saved.GenericOpenAPIError); ok {
		var body map[string]interface{}
		if json.Unmarshal(apiErr.Body(), &body) == nil {
			msg := "API Error"
			if m, ok := body["message"].(string); ok {
				msg = m
			}
			return &APIError{
				Message: msg,
				Details: body,
			}
		}

		if len(apiErr.Body()) > 0 {
			return &APIError{Message: string(apiErr.Body())}
		}

		if apiErr.Model() != nil {
			return &APIError{Message: "API Error", Details: apiErr.Model()}
		}
	}

	return err
}
