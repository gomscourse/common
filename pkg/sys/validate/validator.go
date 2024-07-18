package validate

import "context"

type Condition func(ctx context.Context) error

func Validate(ctx context.Context, conditions ...Condition) error {
	ve := NewValidationErrors()

	for _, c := range conditions {
		if err := c(ctx); err != nil {
			if IsValidationError(err) {
				ve.addError(err.Error())
				continue
			}

			return err
		}
	}

	if ve.Messages == nil {
		return nil
	}

	return ve
}
