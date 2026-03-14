package validation

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

var (
	validatorOnce sync.Once
	validatorInst protovalidate.Validator
	validatorErr  error
)

// ValidateProto runs protovalidate and maps violations to a stable
// "field: message" format for API responses.
func ValidateProto(msg proto.Message) error {
	if msg == nil {
		return fmt.Errorf("request is required")
	}

	v, err := getValidator()
	if err != nil {
		return err
	}
	if err := v.Validate(msg); err != nil {
		return mapValidationError(err)
	}
	return nil
}

func getValidator() (protovalidate.Validator, error) {
	validatorOnce.Do(func() {
		validatorInst, validatorErr = protovalidate.New()
	})
	if validatorErr != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", validatorErr)
	}
	return validatorInst, nil
}

func mapValidationError(err error) error {
	var valErr *protovalidate.ValidationError
	if !errors.As(err, &valErr) {
		return err
	}

	if len(valErr.Violations) == 0 {
		return fmt.Errorf("validation failed")
	}

	messages := make([]string, 0, len(valErr.Violations))
	for _, violation := range valErr.Violations {
		if violation == nil || violation.Proto == nil {
			continue
		}

		fieldPath := protovalidate.FieldPathString(violation.Proto.GetField())
		ruleMessage := strings.TrimSpace(violation.Proto.GetMessage())

		switch {
		case fieldPath == "" && ruleMessage == "":
			continue
		case fieldPath == "":
			messages = append(messages, ruleMessage)
		case ruleMessage == "":
			messages = append(messages, fieldPath+": invalid")
		default:
			messages = append(messages, fieldPath+": "+ruleMessage)
		}
	}

	if len(messages) == 0 {
		return fmt.Errorf("validation failed")
	}
	return errors.New(strings.Join(messages, "; "))
}
