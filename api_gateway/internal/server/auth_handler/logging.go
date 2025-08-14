package handler

import "time"

type AuthLogOption func(*AuthLog)

func WithAuthEmail(email string) AuthLogOption {
	return func(l *AuthLog) { l.Email = email }
}
func WithAuthUserID(userID int64) AuthLogOption {
	return func(l *AuthLog) { l.UserID = userID }
}
func WithAuthError(err string) AuthLogOption {
	return func(l *AuthLog) { l.Error = err }
}
func WithAuthDetails(details string) AuthLogOption {
	return func(l *AuthLog) { l.Details = details }
}

func NewAuthLog(level, event string, opts ...AuthLogOption) AuthLog {
	log := AuthLog{
		Level: level,
		Event: event,
		Time:  time.Now().Format(time.RFC3339),
	}
	for _, opt := range opts {
		opt(&log)
	}
	return log
}
