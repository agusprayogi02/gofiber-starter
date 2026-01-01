package user

// Service defines the interface for authentication service operations
type Service interface {
	Register(user *RegisterRequest) error
	Login(req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error)
	RefreshToken(req *RefreshTokenRequest, ipAddress, userAgent string) (*RefreshTokenResponse, error)
	Logout(refreshToken string) error
	LogoutAll(userID uint) error
	ForgotPassword(req *ForgotPasswordRequest) error
	ResetPassword(req *ResetPasswordRequest) error
	ChangePassword(userID uint, req *ChangePasswordRequest) error
	VerifyEmail(req *VerifyEmailRequest) error
	ResendVerificationEmail(email string) error
	GetActiveSessions(userID uint) ([]SessionResponse, error)
	RevokeSession(sessionID, userID uint) error
}
