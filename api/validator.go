package api

import (
	"github.com/go-playground/validator/v10"
	"simple_bank.sqlc.dev/app/util"
)

// define a custom validator using with validator package for any target fields
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
