package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// check whether currency is supported or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
