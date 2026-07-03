package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Role        string  `json:"role"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	LastLoginAt *string `json:"last_login_at,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=admin viewer"`
}

type BootstrapRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type StatsOverview struct {
	TotalApps     int     `json:"total_apps"`
	ActiveApps    int     `json:"active_apps"`
	TotalUsers    int     `json:"total_users"`
	TotalEmails   int     `json:"total_emails"`
	SentEmails    int     `json:"sent_emails"`
	FailedEmails  int     `json:"failed_emails"`
	SentToday     int     `json:"sent_today"`
	SentThisMonth int     `json:"sent_this_month"`
	SuccessRate   float64 `json:"success_rate"`
}

type AppStats struct {
	AppID         string  `json:"app_id"`
	AppName       string  `json:"app_name"`
	Total         int     `json:"total"`
	Sent          int     `json:"sent"`
	Failed        int     `json:"failed"`
	CreatedByName *string `json:"created_by,omitempty"`
}

type StatsResponse struct {
	Overview StatsOverview `json:"overview"`
	ByApp    []AppStats    `json:"by_app"`
}
