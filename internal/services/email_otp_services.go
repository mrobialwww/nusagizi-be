package services

import (
	"context"
	"errors"

	"nusagizi_be/internal/auth/auth0client"
)

type EmailOTPService struct {
	verifier   *auth0client.IDTokenVerifier
	management *auth0client.ManagementClient
}

func NewEmailOTPService(verifier *auth0client.IDTokenVerifier, management *auth0client.ManagementClient) *EmailOTPService {
	return &EmailOTPService{verifier: verifier, management: management}
}

func (s *EmailOTPService) ConfirmEmail(ctx context.Context, idToken string) error {
	claims, err := s.verifier.Verify(idToken)
	if err != nil {
		return errors.New("id_token tidak valid")
	}

	// Auth0 sets email_verified=true on this Passwordless "email" connection
	// token upon successful OTP entry, proving email ownership.
	if !claims.EmailVerified {
		return errors.New("email belum terverifikasi oleh Auth0")
	}

	if claims.Email == "" {
		return errors.New("id_token tidak memiliki klaim email")
	}

	user, err := s.management.FindDatabaseUserByEmail(ctx, claims.Email)
	if err != nil {
		return errors.New("akun database untuk email ini tidak ditemukan")
	}

	if err := s.management.MarkEmailVerified(ctx, user.UserID); err != nil {
		return errors.New("gagal memperbarui status verifikasi email")
	}

	return nil
}
