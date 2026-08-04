package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"nusagizi_be/internal/models/checkin"
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
func (s *CheckinService) ValidateToken(ctx context.Context, token string, userID string) (*checkin.ValidateResult, error) {
	// 1. Check if token exists
	data, err := s.repo.GetToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: token not found", repository.ErrNotFound)
		}
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(data.ExpiresAt) {
		return nil, fmt.Errorf("%w: token expired", repository.ErrForbidden)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// 3. Record check-in log
	err = s.repo.SaveLog(ctx, token, data.ChildID, userUUID)
	if err != nil {
		return nil, err
	}

	// 4. Find caregiver profile belonging to this user
	caregiverProfileID, err := s.caregiverRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Reject if user doesn't have a caregiver profile
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: no caregiver profile associated with this account", repository.ErrForbidden)
		}
		return nil, err
	}

	// 5. Add Caregiver Engagement
	var conflictErr error
	engagementID, err := s.repo.CreateCaregiverEngagement(ctx, data.ChildID, caregiverProfileID)
	if err != nil {
		// If fails due to conflict (already engaged), we still proceed to return the enriched data
		if errors.Is(err, repository.ErrConflict) || err.Error() == "active engagement already exists" {
			conflictErr = repository.ErrConflict
		} else {
			return nil, err
		}
	}

	// 6. Fetch additional child and mother information
	childInfo, err := s.repo.GetChildInfo(ctx, data.ChildID)
	if err != nil {
		return nil, err
	}

	var engagementStr string
	if engagementID != uuid.Nil {
		engagementStr = engagementID.String()
	}

	result := &checkin.ValidateResult{
		Valid:        true,
		EngagementID: engagementStr,
		CheckedInAt:  time.Now(),
		ChildName:    childInfo.ChildName,
		ChildAge:     formatAge(childInfo.BirthDate),
		MotherName:   childInfo.MotherName,
	}

	return result, conflictErr
}

// formatAge calculates the age in "X tahun Y bulan" format.
func formatAge(birthDate time.Time) string {
	now := time.Now()
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	if now.Day() < birthDate.Day() {
		months--
	}
	if months < 0 {
		years--
		months += 12
	}
	if years > 0 && months > 0 {
		return fmt.Sprintf("%d tahun %d bulan", years, months)
	} else if years > 0 {
		return fmt.Sprintf("%d tahun", years)
	}
	return fmt.Sprintf("%d bulan", months)
}
