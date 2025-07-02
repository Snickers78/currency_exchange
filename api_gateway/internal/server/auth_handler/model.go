package handler

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthLog struct {
	Level   string `json:"level"`
	Event   string `json:"event"`
	Email   string `json:"email,omitempty"`
	UserID  int64  `json:"user_id,omitempty"`
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
	Time    string `json:"time"`
}
