package handler

import "time"

type ExchangeLogOption func(*ExchangeLog)

func WithBaseCurrency(base string) ExchangeLogOption {
	return func(l *ExchangeLog) { l.BaseCurrency = base }
}
func WithTargetCurrency(target string) ExchangeLogOption {
	return func(l *ExchangeLog) { l.TargetCurrency = target }
}
func WithAmount(amount float64) ExchangeLogOption {
	return func(l *ExchangeLog) { l.Amount = amount }
}
func WithRate(rate float64) ExchangeLogOption {
	return func(l *ExchangeLog) { l.Rate = rate }
}
func WithCurrencyName(name string) ExchangeLogOption {
	return func(l *ExchangeLog) { l.CurrencyName = name }
}
func WithExchangeError(err string) ExchangeLogOption {
	return func(l *ExchangeLog) { l.Error = err }
}
func WithExchangeDetails(details string) ExchangeLogOption {
	return func(l *ExchangeLog) { l.Details = details }
}

func NewExchangeLog(level, event string, opts ...ExchangeLogOption) ExchangeLog {
	log := ExchangeLog{
		Level: level,
		Event: event,
		Time:  time.Now().Format(time.RFC3339),
	}
	for _, opt := range opts {
		opt(&log)
	}
	return log
}
