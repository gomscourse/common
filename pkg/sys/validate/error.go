package validate

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type ValidationErrors struct {
	Messages []string `json:"error_messages"`
}

func (ve *ValidationErrors) addError(message string) {
	ve.Messages = append(ve.Messages, message)
}

func NewValidationErrors(messages ...string) *ValidationErrors {
	return &ValidationErrors{
		Messages: messages,
	}
}

func (ve *ValidationErrors) Error() string {
	data, err := json.Marshal(ve.Messages)
	if err != nil {
		return err.Error()
	}

	return string(data)
}

func IsValidationError(err error) bool {
	var ve *ValidationErrors
	return errors.As(err, &ve)
}
