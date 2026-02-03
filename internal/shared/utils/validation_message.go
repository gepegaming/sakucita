package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GenerateMessageValidation(err error) map[string]any {
	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		messages := make(map[string]any, len(vErr))
		for _, v := range vErr {
			fieldName := v.Field()

			switch v.Tag() {
			case "alpha":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters", fieldName)
			case "alpha_dash":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters, dashes, and underscores", fieldName)
			case "alpha_num":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters and numbers", fieldName)
			case "date":
				messages[fieldName] = fmt.Sprintf("%s is not a valid date, e.g: 2006-01-02", fieldName)
			case "email":
				messages[fieldName] = fmt.Sprintf("%s is not valid email", v.Value())
			case "gt":
				messages[fieldName] = fmt.Sprintf("%s must be greater than %s", fieldName, v.Param())
			case "gte":
				messages[fieldName] = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, v.Param())
			case "hex":
				messages[fieldName] = fmt.Sprintf("%s must be a valid hexadecimal string", fieldName)
			case "ip":
				messages[fieldName] = fmt.Sprintf("%s is not a valid IP address", fieldName)
			case "len":
				messages[fieldName] = fmt.Sprintf("%s must be exactly %s characters", fieldName, v.Param())
			case "lt":
				messages[fieldName] = fmt.Sprintf("%s must be less than %s", fieldName, v.Param())
			case "lte":
				messages[fieldName] = fmt.Sprintf("%s must be less than or equal to %s", fieldName, v.Param())
			case "max":
				messages[fieldName] = fmt.Sprintf("%s must be less than %s characters", fieldName, v.Param())
			case "min":
				messages[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, v.Param())
			case "numeric":
				messages[fieldName] = fmt.Sprintf("%s must be a numeric value", fieldName)
			case "oneof":
				messages[fieldName] = fmt.Sprintf("%s must be one of: [%s]", fieldName, v.Param())
			case "required":
				messages[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "eq":
				messages[fieldName] = fmt.Sprintf("%s must be equal to %s", fieldName, v.Param())
			case "required_if":
				messages[fieldName] = fmt.Sprintf("%s is required when %s is %s", fieldName, v.Field(), v.Param())
			case "timezone":
				messages[fieldName] = fmt.Sprintf("%s is not a valid timezone e.g: UTC,+08:00,Asia,Jakarta,America,New_York", fieldName)
			case "url":
				messages[fieldName] = fmt.Sprintf("%s is not a valid URL", fieldName)
			case "uuid":
				messages[fieldName] = fmt.Sprintf("%s is not a valid UUID", fieldName)
			}
		}
		return messages
	}
	return map[string]any{"error": err.Error()}
}
