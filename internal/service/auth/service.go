package auth

import (
	"time"

	"starter-gofiber/internal/domain/user"
	"starter-gofiber/internal/worker"
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
	// Create email verification token
	verificationToken, err := crypto.GenerateRandomToken()
	if err != nil {
		return &apierror.InternalServerError{
			Message: "Failed to generate verification token",
			Order:   "S-Register-3",
		}
	}

	emailVerification := &user.EmailVerification{
		UserID:    userEntity.ID,
		Token:     verificationToken,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := s.userRepo.CreateEmailVerification(emailVerification); err != nil {
		return &apierror.InternalServerError{
			Message: "Failed to create email verification",
			Order:   "S-Register-4",
		}
	}

	// Send verification email via background worker
	if _, err := worker.EnqueueEmailVerification(userEntity.Email, verificationToken); err != nil {
		// Log error but don't fail registration
		// Email will be sent when worker processes the queue
	}

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

	// Send password reset email via background worker
	if _, err := worker.EnqueueEmailPasswordReset(usr.Email, resetToken); err != nil {
		// Log error but don't fail the request
		// Email will be sent when worker processes the queue
	}

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

	// Send verification email via background worker
	if _, err := worker.EnqueueEmailVerification(usr.Email, verificationToken); err != nil {
		// Log error but don't fail the request
		// Email will be sent when worker processes the queue
	}

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

// Profile operations
func (s *AuthService) GetProfile(userID uint) (*user.GetProfileResponse, error) {
	usr, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	response := user.GetProfileResponse{}.FromEntity(*usr)
	return &response, nil
}

func (s *AuthService) UpdateProfile(userID uint, req *user.UpdateProfileRequest) (*user.GetProfileResponse, error) {
	usr, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Update fields if provided
	if req.Name != "" {
		usr.Name = req.Name
	}
	if req.Bio != "" {
		usr.Bio = req.Bio
	}

	if err := s.userRepo.Update(usr); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	response := user.GetProfileResponse{}.FromEntity(*usr)
	return &response, nil
}

func (s *AuthService) UpdateAvatar(userID uint, avatarPath string) (*user.GetProfileResponse, error) {
	usr, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Delete old avatar if exists (optional - can be implemented later)
	// storage.DeleteFile(&usr.Avatar, variables.AVATAR_PATH)

	// Update avatar path
	usr.Avatar = avatarPath

	if err := s.userRepo.Update(usr); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	response := user.GetProfileResponse{}.FromEntity(*usr)
	return &response, nil
}

// Preferences operations
func (s *AuthService) GetPreferences(userID uint) (*user.GetPreferencesResponse, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Get or create preferences
	prefs, err := s.userRepo.FindPreferencesByUserID(userID)
	if err != nil {
		// Preferences not found, create default
		defaultPrefs := &user.UserPreferences{
			UserID: userID,
		}
		defaultData, _ := defaultPrefs.GetData()
		defaultPrefs.SetData(defaultData)

		if err := s.userRepo.CreatePreferences(defaultPrefs); err != nil {
			return nil, &apierror.InternalServerError{
				Message: err.Error(),
				Order:   "S2",
			}
		}
		prefs = defaultPrefs
	}

	data, err := prefs.GetData()
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	return &user.GetPreferencesResponse{
		Preferences: data,
		UpdatedAt:   prefs.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) UpdatePreferences(userID uint, req *user.UpdatePreferencesRequest) (*user.GetPreferencesResponse, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, &apierror.NotFoundError{
			Message: "User not found",
			Order:   "S1",
		}
	}

	// Get or create preferences
	prefs, err := s.userRepo.FindPreferencesByUserID(userID)
	if err != nil {
		// Preferences not found, create new
		prefs = &user.UserPreferences{
			UserID: userID,
		}
	}

	// Get current preferences data
	data, err := prefs.GetData()
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S2",
		}
	}

	// Update fields if provided
	if req.EmailNotifications != nil {
		data.EmailNotifications = *req.EmailNotifications
	}
	if req.PushNotifications != nil {
		data.PushNotifications = *req.PushNotifications
	}
	if req.SMSNotifications != nil {
		data.SMSNotifications = *req.SMSNotifications
	}
	if req.ProfileVisibility != nil {
		data.ProfileVisibility = *req.ProfileVisibility
	}
	if req.ShowEmail != nil {
		data.ShowEmail = *req.ShowEmail
	}
	if req.ShowOnlineStatus != nil {
		data.ShowOnlineStatus = *req.ShowOnlineStatus
	}
	if req.Theme != nil {
		data.Theme = *req.Theme
	}
	if req.Language != nil {
		data.Language = *req.Language
	}
	if req.Timezone != nil {
		data.Timezone = *req.Timezone
	}
	if req.Custom != nil {
		if data.Custom == nil {
			data.Custom = make(map[string]interface{})
		}
		for k, v := range req.Custom {
			data.Custom[k] = v
		}
	}

	// Save preferences
	if err := prefs.SetData(data); err != nil {
		return nil, &apierror.InternalServerError{
			Message: err.Error(),
			Order:   "S3",
		}
	}

	// Create or update preferences
	if prefs.ID == 0 {
		if err := s.userRepo.CreatePreferences(prefs); err != nil {
			return nil, &apierror.InternalServerError{
				Message: err.Error(),
				Order:   "S4",
			}
		}
	} else {
		if err := s.userRepo.UpdatePreferences(prefs); err != nil {
			return nil, &apierror.InternalServerError{
				Message: err.Error(),
				Order:   "S4",
			}
		}
	}

	return &user.GetPreferencesResponse{
		Preferences: data,
		UpdatedAt:   prefs.UpdatedAt.Format(time.RFC3339),
	}, nil
}
