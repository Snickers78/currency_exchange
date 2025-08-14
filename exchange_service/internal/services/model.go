package currency

type exchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"conversion_rates"`
	ErrorType          string             `json:"error-type,omitempty"`
}

type exchangeCurrencyAPIResponse struct {
	Result             string  `json:"result"`
	Documentation      string  `json:"documentation"`
	TermsOfUse         string  `json:"terms_of_use"`
	TimeLastUpdateUnix int64   `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string  `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64   `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string  `json:"time_next_update_utc"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
	ConversionResult   float64 `json:"conversion_result"`
	ErrorType          string  `json:"error-type,omitempty"`
}
