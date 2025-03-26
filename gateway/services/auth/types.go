package auth

type SignUpInput struct {
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type AuthResponse struct {
	User  UserResponse  `json:"user"`
	Token *TokenResponse `json:"token"`
}
