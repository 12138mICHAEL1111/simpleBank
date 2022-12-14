package api

import (
	"github.com/12138mICHAEL1111/simplebank/util"
	"github.com/go-playground/validator/v10"
)
//验证currency
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency,ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}