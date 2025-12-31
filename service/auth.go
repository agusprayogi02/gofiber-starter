package service

import (
	"time"

	"starter-gofiber/dto"
	"starter-gofiber/entity"
	"starter-gofiber/helper"
	"starter-gofiber/repository"
)

type AuthService struct {
	userR          *repository.UserRepository
	refreshTokenR  *repository.RefreshTokenRepository
	passwordResetR *repository.PasswordResetRepository
	emailVerifR    *repository.EmailVerificationRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	passwordResetRepo *repository.PasswordResetRepository,
	emailVerifRepo *repository.EmailVerificationRepository,
) *AuthService {
	return &AuthService{
		userR:          userRepo,
		refreshTokenR:  refreshTokenRepo,
		passwordResetR: passwordResetRepo,
		emailVerifR:    emailVerifRepo,
	}
}

func (s *AuthService) Register(user *dto.RegisterRequest) error {
	if err := s.userR.ExistEmail(user.Email); err == nil {
		return &helper.BadRequestError{Message: "Email already exists", Order: "S1"}
	}

	password, err := helper.HashPassword(user.Password)
	if err != nil {
		return &helper.BadRequestError{Message: "Failed to hash password", Order: "S2"}
	}

	userEntity := user.ToEntity()
	userEntity.Password = password

	err = s.userR.Create(userEntity)
	if err != nil {
		return &helper.InternalServerError{Message: err.Error(), Order: "S3"}
	}

	// Create email verification token (optional, untuk kirim email nanti)
	// Uncomment jika sudah ada email service
	// verificationToken, _ := helper.GenerateRandomToken()
	// emailVerification := &entity.EmailVerification{
	// 	UserID:    userEntity.ID,
	// 	Token:     verificationToken,
	// 	ExpiresAt: time.Now().Add(time.Hour * 24),
	// }
	// s.emailVerifR.Create(emailVerification)
	// TODO: Send verification email

	return nil
}

func (s *AuthService) Login(req *dto.LoginRequest, ipAddress, userAgent string) (resp *dto.LoginResponse, err error) {
	user, err := s.userR.FindByEmail(req.Email)
	if err != nil {
		// Record failed login attempt
		helper.RecordAuthAttempt(false)
		return nil, &helper.UnauthorizedError{
			Message: "Email not registered!",
			Order:   "S1",
		}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		// Record failed login attempt
		helper.RecordAuthAttempt(false)
		return nil, &helper.UnauthorizedError{
			Message: "Password is wrong!",
			Order:   "S2",
		}
	}

	userClaims := dto.UserClaims{}.FromEntity(*user)
	token, err := helper.GenerateJWT(userClaims)
	if err != nil {
		helper.RecordAuthAttempt(false)
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	refreshToken, err := helper.GenerateRefreshToken(userClaims)
	if err != nil {
		helper.RecordAuthAttempt(false)
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Save refresh token to database
	refreshTokenEntity := &entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.refreshTokenR.Create(refreshTokenEntity); err != nil {
		helper.RecordAuthAttempt(false)
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S5",
		}
	}

	// Record successful login
	helper.RecordAuthAttempt(true)

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         dto.UserResponse{}.FromEntity(*user),
	}, nil
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest, ipAddress, userAgent string) (*dto.RefreshTokenResponse, error) {
	// Validate refresh token
	tokenEntity, err := s.refreshTokenR.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, &helper.UnauthorizedError{
			Message: "Invalid refresh token",
			Order:   "S1",
		}
	}

	// Generate new tokens
	userClaims := dto.UserClaims{}.FromEntity(tokenEntity.User)
	newToken, err := helper.GenerateJWT(userClaims)
	if err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	newRefreshToken, err := helper.GenerateRefreshToken(userClaims)
	if err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Revoke old refresh token
	if err := s.refreshTokenR.RevokeToken(req.RefreshToken); err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Save new refresh token
	newRefreshTokenEntity := &entity.RefreshToken{
		UserID:    tokenEntity.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.refreshTokenR.Create(newRefreshTokenEntity); err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S5",
		}
	}

	return &dto.RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	if err := s.refreshTokenR.RevokeToken(refreshToken); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}

func (s *AuthService) LogoutAll(userID uint) error {
	if err := s.refreshTokenR.RevokeAllUserTokens(userID); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}

func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userR.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate reset token
	resetToken, err := helper.GenerateRandomToken()
	if err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	// Delete old reset tokens for this user
	s.passwordResetR.DeleteUserResets(user.ID)

	// Create new reset token
	passwordReset := &entity.PasswordReset{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Hour * 1), // 1 hour
	}
	if err := s.passwordResetR.Create(passwordReset); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// TODO: Send email with reset token
	// sendPasswordResetEmail(user.Email, resetToken)

	return nil
}

func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// Find reset token
	resetToken, err := s.passwordResetR.FindByToken(req.Token)
	if err != nil {
		return &helper.BadRequestError{
			Message: "Invalid or expired reset token",
			Order:   "S1",
		}
	}

	// Hash new password
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// Update user password
	user := resetToken.User
	user.Password = hashedPassword
	if err := s.userR.Update(&user); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Mark token as used
	if err := s.passwordResetR.MarkAsUsed(req.Token); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Revoke all user sessions for security
	s.refreshTokenR.RevokeAllUserTokens(user.ID)

	return nil
}

func (s *AuthService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userR.FindByID(userID)
	if err != nil {
		return &helper.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Verify old password
	if err := helper.VerifyPassword(user.Password, req.OldPassword); err != nil {
		return &helper.BadRequestError{
			Message: "Old password is incorrect",
			Order:   "S2",
		}
	}

	// Hash new password
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userR.Update(user); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// Revoke all sessions except current
	s.refreshTokenR.RevokeAllUserTokens(userID)

	return nil
}

func (s *AuthService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	// Find verification token
	verification, err := s.emailVerifR.FindByToken(req.Token)
	if err != nil {
		return &helper.BadRequestError{
			Message: "Invalid or expired verification token",
			Order:   "S1",
		}
	}

	// Update user email verified status
	user := verification.User
	user.EmailVerified = true
	if err := s.userR.Update(&user); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// Mark verification as used
	if err := s.emailVerifR.MarkAsVerified(req.Token); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	return nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.userR.FindByEmail(email)
	if err != nil {
		return &helper.NotFoundError{
			Message: "Email not found",
			Order:   "S1",
		}
	}

	if user.EmailVerified {
		return &helper.BadRequestError{
			Message: "Email already verified",
			Order:   "S2",
		}
	}

	// Generate new verification token
	verificationToken, err := helper.GenerateRandomToken()
	if err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Delete old verification tokens
	s.emailVerifR.DeleteUserVerifications(user.ID)

	// Create new verification
	emailVerification := &entity.EmailVerification{
		UserID:    user.ID,
		Token:     verificationToken,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := s.emailVerifR.Create(emailVerification); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S4",
		}
	}

	// TODO: Send verification email
	// sendVerificationEmail(user.Email, verificationToken)

	return nil
}

func (s *AuthService) GetActiveSessions(userID uint) ([]dto.SessionResponse, error) {
	sessions, err := s.refreshTokenR.GetUserActiveSessions(userID)
	if err != nil {
		return nil, &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}

	var response []dto.SessionResponse
	for _, session := range sessions {
		response = append(response, dto.SessionResponse{}.FromEntity(session))
	}

	return response, nil
}

func (s *AuthService) RevokeSession(sessionID, userID uint) error {
	if err := s.refreshTokenR.RevokeSessionByID(sessionID, userID); err != nil {
		return &helper.InternalServerError{
			Message: err.Error(),
			Order:   "S1",
		}
	}
	return nil
}
