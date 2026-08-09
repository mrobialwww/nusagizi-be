package services

import (
	"context"
	"fmt"

	"nusagizi_be/internal/models/caregiver"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

type CaregiverService struct {
	repo       *repository.CaregiverRepository
	motherRepo *repository.MotherProfileRepository
}

func NewCaregiverService(
	repo *repository.CaregiverRepository,
	motherRepo *repository.MotherProfileRepository,
) *CaregiverService {
	return &CaregiverService{
		repo:       repo,
		motherRepo: motherRepo,
	}
}

// GetCaregiverEngagements (endpoint: 52)
func (s *CaregiverService) GetCaregiverEngagements(ctx context.Context, userID string) ([]caregiver.CaregiverEngagementResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetCaregiverEngagements(ctx, motherProfileID, false)
}

// GetRevokedEngagements (endpoint: 53)
func (s *CaregiverService) GetRevokedEngagements(ctx context.Context, userID string) ([]caregiver.CaregiverEngagementResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetCaregiverEngagements(ctx, motherProfileID, true)
}

// DeleteCaregiverEngagement (endpoint: 54)
func (s *CaregiverService) DeleteCaregiverEngagement(ctx context.Context, userID string, engagementID uuid.UUID) error {
	// Must verify if engagement belongs to this mother's children
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: mother profile not found", repository.ErrForbidden)
	}

	// We can check if engagement belongs to mother by checking the list
	engagements, err := s.repo.GetCaregiverEngagements(ctx, motherProfileID, false)
	if err != nil {
		return err
	}

	owns := false
	for _, ce := range engagements {
		if ce.CaregiverEngagementID == engagementID {
			owns = true
			break
		}
	}
	if !owns {
		return fmt.Errorf("%w: engagement does not belong to you", repository.ErrForbidden)
	}

	return s.repo.DeleteCaregiverEngagement(ctx, engagementID)
}

// GetCaregiverChildren (endpoint: 56)
func (s *CaregiverService) GetCaregiverChildren(ctx context.Context, userID string) ([]caregiver.ChildSimpleResponse, error) {
	// Automatically resolve caregiver profile from the logged-in user
	caregiverProfileID, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: caregiver profile not found", repository.ErrForbidden)
	}

	return s.repo.GetCaregiverChildren(ctx, caregiverProfileID)
}
