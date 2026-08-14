package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"nusagizi_be/internal/repository"
	"time"

	"github.com/google/uuid"
)

const qrDuration = 5 * time.Minute

type CheckinService struct {
	repo          *repository.CheckinRepository
	caregiverRepo *repository.CaregiverRepository
}

func NewCheckinService(repo *repository.CheckinRepository, caregiverRepo *repository.CaregiverRepository) *CheckinService {
	return &CheckinService{repo: repo, caregiverRepo: caregiverRepo}
}

// GenerateToken (Endpoint: 65)
func (s *CheckinService) GenerateToken(ctx context.Context, childID uuid.UUID) (string, time.Time, error) {
	// Jalankan lazy cleanup secara asynchronous di background
	go func() {
		cleanupCtx := context.Background()
		_ = s.repo.DeleteExpiredTokens(cleanupCtx)
	}()

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(qrDuration)

	err := s.repo.SaveToken(ctx, token, childID, expiresAt)
	return token, expiresAt, err
}

// ValidateToken (Endpoint: 66)
func (s *CheckinService) ValidateToken(ctx context.Context, token string, userID string) (uuid.UUID, error) {
	// 1. Check if token exists
	data, err := s.repo.GetToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return uuid.Nil, fmt.Errorf("%w: token not found", repository.ErrNotFound)
		}
		return uuid.Nil, err
	}

	// Check if token is expired
	if time.Now().After(data.ExpiresAt) {
		return uuid.Nil, fmt.Errorf("%w: token expired", repository.ErrForbidden)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}

	// 3. Record check-in log
	err = s.repo.SaveLog(ctx, token, data.ChildID, userUUID)
	if err != nil {
		return uuid.Nil, err
	}

	// 4. Find caregiver profile belonging to this user
	caregiverProfileID, err := s.caregiverRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Reject if user doesn't have a caregiver profile
		if errors.Is(err, repository.ErrNotFound) {
			return uuid.Nil, fmt.Errorf("%w: no caregiver profile associated with this account", repository.ErrForbidden)
		}
		return uuid.Nil, err
	}

	// 5. Add Caregiver Engagement
	engagementID, err := s.repo.CreateCaregiverEngagement(ctx, data.ChildID, caregiverProfileID)
	if err != nil {
		// If fails due to conflict (already engaged), return the existing ID and ErrConflict
		if errors.Is(err, repository.ErrConflict) {
			return engagementID, repository.ErrConflict
		}
		return uuid.Nil, err
	}

	return engagementID, nil
}
