package services

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

type ChildService struct {
	repo          *repository.ChildRepository
	motherRepo    *repository.MotherProfileRepository
	caregiverRepo *repository.CaregiverRepository
}

func NewChildService(repo *repository.ChildRepository, motherRepo *repository.MotherProfileRepository, caregiverRepo *repository.CaregiverRepository) *ChildService {
	return &ChildService{repo: repo, motherRepo: motherRepo, caregiverRepo: caregiverRepo}
}

// CreateChild creates a new child profile within a DB transaction (Endpoint: 7)
func (s *ChildService) CreateChild(ctx context.Context, userID string, input *models.CreateChildInput) (uuid.UUID, error) {
	// retrieve mother_profile_id from user_id
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return uuid.Nil, err
	}

	return s.repo.Create(ctx, motherProfileID, input)
}

// UpdateChild updates a child profile (Endpoint: 8)
func (s *ChildService) UpdateChild(ctx context.Context, userID string, childID uuid.UUID, input *models.UpdateChildInput) error {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.repo); err != nil {
		return err
	}
	return s.repo.Update(ctx, childID, input)
}

// DeleteChild soft-deletes a child (Endpoint: 9)
func (s *ChildService) DeleteChild(ctx context.Context, userID string, childID uuid.UUID) error {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.repo); err != nil {
		return err
	}

	return s.repo.SoftDelete(ctx, childID)
}

// GetChildDetail returns the full child detail (Endpoint: 10)
func (s *ChildService) GetChildDetail(ctx context.Context, userID string, childID uuid.UUID) (*models.ChildDetailResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.repo); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, childID)
}

// GetChildSimple returns the simple child detail (Endpoint: 12 - Get profil anak versi ringan)
func (s *ChildService) GetChildSimple(ctx context.Context, userID string, childID uuid.UUID) (*models.ChildSimpleResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.repo); err != nil {
		return nil, err
	}

	return s.repo.GetSimpleByID(ctx, childID)
}

// GetChildrenByMother returns the lightweight children list (Endpoint: 11)
func (s *ChildService) GetChildrenByMother(ctx context.Context, userID string) ([]models.ChildListItem, error) {
	// retrieve mother_profile_id from user_id
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return nil, err
	}

	return s.repo.GetByMotherProfileID(ctx, motherProfileID)
}
