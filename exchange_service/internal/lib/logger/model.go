package logg

type ExchangeLog struct {
	Level          string  `json:"level"`
	Event          string  `json:"event"`
	BaseCurrency   string  `json:"base_currency,omitempty"`
	TargetCurrency string  `json:"target_currency,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	Rate           float64 `json:"rate,omitempty"`
	CurrencyName   string  `json:"currency_name,omitempty"`
	Error          string  `json:"error,omitempty"`
	Details        string  `json:"details,omitempty"`
	Time           string  `json:"time"`
}
