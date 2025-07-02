package handler

type rateRequest struct {
	BaseCurrency   string `json:"base_currency" binding:"required"`
	TargetCurrency string `json:"target_currency" binding:"required"`
}

type exchangeRequest struct {
	BaseCurrency   string  `json:"base_currency" binding:"required"`
	TargetCurrency string  `json:"target_currency" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
}

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
