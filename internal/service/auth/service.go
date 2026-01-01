package auth

import (
	"time"

	"starter-gofiber/internal/domain/user"
	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/crypto"
)

type AuthService struct {
	userRepo user.Repository
}

func NewAuthService(userRepo user.Repository) user.Service {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(req *user.RegisterRequest) error {
	if err := s.userRepo.ExistEmail(req.Email); err == nil {
		return &apierror.BadRequestError{Message: "Email already exists", Order: "S1"}
	}

	password, err := crypto.HashPassword(req.Password)
	if err != nil {
		return &apierror.BadRequestError{Message: "Failed to hash password", Order: "S2"}
	}

	userEntity := req.ToEntity()
	userEntity.Password = password

	err = s.userRepo.Create(&userEntity)
	if err != nil {
		return &apierror.InternalServerError{Message: err.Error(), Order: "S3"}
	}

	// Create email verification token (optional, untuk kirim email nanti)
	// Uncomment jika sudah ada email service
	// verificationToken, _ := crypto.GenerateRandomToken()
	// emailVerification := &user.EmailVerification{
	// 	UserID:    userEntity.ID,
	// 	Token:     verificationToken,
	// 	ExpiresAt: time.Now().Add(time.Hour * 24),
	// }
	// s.userRepo.CreateEmailVerification(emailVerification)
	// TODO: Send verification email

	return nil
}

func (s *AuthService) Login(req *user.LoginRequest, ipAddress, userAgent string) (resp *user.LoginResponse, err error) {
	usr, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Record failed login attempt
		return nil, &apierror.UnauthorizedError{
			Message: "Email not registered!",
			Order:   "S1",
		}
	}

	if err := crypto.VerifyPassword(usr.Password, req.Password); err != nil {
		// Record failed login attempt
		return nil, &apierror.UnauthorizedError{
			Message: "Password is wrong!",
			Order:   "S2",
		}
	}

	userClaims := user.UserClaims{}.FromEntity(*usr)
	token, err := crypto.GenerateJWT(userClaims)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	refreshToken, err := crypto.GenerateRefreshToken(userClaims)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Save refresh token to database
	refreshTokenEntity := &user.RefreshToken{
		UserID:    usr.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.userRepo.CreateRefreshToken(refreshTokenEntity); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S5",
		}
	}

	// Record successful login

	return &user.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user.UserResponse{}.FromEntity(*usr),
	}, nil
}

func (s *AuthService) RefreshToken(req *user.RefreshTokenRequest, ipAddress, userAgent string) (*user.RefreshTokenResponse, error) {
	// Validate refresh token
	tokenEntity, err := s.userRepo.FindRefreshTokenByToken(req.RefreshToken)
	if err != nil {
		return nil, &apierror.UnauthorizedError{
			Message: "Invalid refresh token",
			Order:   "S1",
		}
	}

	// Generate new tokens
	userClaims := user.UserClaims{}.FromEntity(tokenEntity.User)
	newToken, err := crypto.GenerateJWT(userClaims)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	newRefreshToken, err := crypto.GenerateRefreshToken(userClaims)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Revoke old refresh token
	if err := s.userRepo.RevokeRefreshToken(req.RefreshToken); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Save new refresh token
	newRefreshTokenEntity := &user.RefreshToken{
		UserID:    tokenEntity.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.userRepo.CreateRefreshToken(newRefreshTokenEntity); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S5",
		}
	}

	return &user.RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	if err := s.userRepo.RevokeRefreshToken(refreshToken); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}

func (s *AuthService) LogoutAll(userID uint) error {
	if err := s.userRepo.RevokeAllUserTokens(userID); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}

func (s *AuthService) ForgotPassword(req *user.ForgotPasswordRequest) error {
	usr, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	resetToken, err := crypto.GenerateRandomToken()
	if err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	// Create new reset token
	passwordReset := &user.PasswordReset{
		UserID:    usr.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Hour * 1), // 1 hour
	}
	if err := s.userRepo.CreatePasswordReset(passwordReset); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// TODO: Send email with reset token
	// sendPasswordResetEmail(usr.Email, resetToken)

	return nil
}

func (s *AuthService) ResetPassword(req *user.ResetPasswordRequest) error {
	// Find reset token
	resetToken, err := s.userRepo.FindPasswordResetByToken(req.Token)
	if err != nil {
		return &apierror.BadRequestError{
			Message: "Invalid or expired reset token",
			Order:   "S1",
		}
	}

	// Hash new password
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// Update user password
	usr := resetToken.User
	usr.Password = hashedPassword
	if err := s.userRepo.Update(&usr); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Mark token as used
	if err := s.userRepo.MarkPasswordResetAsUsed(req.Token); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Revoke all user sessions for security
	s.userRepo.RevokeAllUserTokens(usr.ID)

	return nil
}

func (s *AuthService) ChangePassword(userID uint, req *user.ChangePasswordRequest) error {
	// Get user
	usr, err := s.userRepo.FindByID(userID)
	if err != nil {
		return &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Verify old password
	if err := crypto.VerifyPassword(usr.Password, req.OldPassword); err != nil {
		return &apierror.BadRequestError{
			Message: "Old password is incorrect",
			Order:   "S2",
		}
	}

	// Hash new password
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Update password
	usr.Password = hashedPassword
	if err := s.userRepo.Update(usr); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Revoke all sessions except current
	s.userRepo.RevokeAllUserTokens(userID)

	return nil
}

func (s *AuthService) VerifyEmail(req *user.VerifyEmailRequest) error {
	// Find verification token
	verification, err := s.userRepo.FindEmailVerificationByToken(req.Token)
	if err != nil {
		return &apierror.BadRequestError{
			Message: "Invalid or expired verification token",
			Order:   "S1",
		}
	}

	// Update user email verified status
	usr := verification.User
	usr.EmailVerified = true
	if err := s.userRepo.Update(&usr); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// Mark verification as used
	if err := s.userRepo.MarkEmailAsVerified(usr.ID); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	return nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	usr, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return &apierror.NotFoundError{
			Message: "Email not found",
			Order:   "S1",
		}
	}

	if usr.EmailVerified {
		return &apierror.BadRequestError{
			Message: "Email already verified",
			Order:   "S2",
		}
	}

	// Generate new verification token
	verificationToken, err := crypto.GenerateRandomToken()
	if err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Create new verification
	emailVerification := &user.EmailVerification{
		UserID:    usr.ID,
		Token:     verificationToken,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := s.userRepo.CreateEmailVerification(emailVerification); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// TODO: Send verification email
	// sendVerificationEmail(usr.Email, verificationToken)

	return nil
}

func (s *AuthService) GetActiveSessions(userID uint) ([]user.SessionResponse, error) {
	sessions, err := s.userRepo.FindUserRefreshTokens(userID)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	var response []user.SessionResponse
	for _, session := range sessions {
		response = append(response, user.SessionResponse{}.FromEntity(session))
	}

	return response, nil
}

func (s *AuthService) RevokeSession(sessionID, userID uint) error {
	if err := s.userRepo.RevokeSessionByID(sessionID, userID); err != nil {
		return &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}
