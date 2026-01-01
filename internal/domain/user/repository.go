package user

// Repository defines the interface for user repository operations
type Repository interface {
	// User operations
	Create(user *User) error
	ExistEmail(email string) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error

	// RefreshToken operations
	CreateRefreshToken(token *RefreshToken) error
	FindRefreshTokenByToken(token string) (*RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID uint) error
	FindUserRefreshTokens(userID uint) ([]RefreshToken, error)
	RevokeSessionByID(sessionID uint, userID uint) error

	// EmailVerification operations
	CreateEmailVerification(verification *EmailVerification) error
	FindEmailVerificationByToken(token string) (*EmailVerification, error)
	MarkEmailAsVerified(userID uint) error

	// PasswordReset operations
	CreatePasswordReset(reset *PasswordReset) error
	FindPasswordResetByToken(token string) (*PasswordReset, error)
	MarkPasswordResetAsUsed(token string) error

	// APIKey operations
	CreateAPIKey(apiKey *APIKey) error
	FindAPIKeyByHash(hash string) (*APIKey, error)
	UpdateAPIKey(apiKey *APIKey) error
}
