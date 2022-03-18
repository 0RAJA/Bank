package valid

import "github.com/go-playground/validator/v10"

// 创建对于货币的校验
var (
	currencys = []string{"USD", "EUR", "RMB"}
)

func isSupportedCurrencies(currency string) bool {
	for _, s := range currencys {
		if currency == s {
			return true
		}
	}
	return false
}

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return isSupportedCurrencies(currency)
	}
	return false
}
