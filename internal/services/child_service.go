package services

import (
	"context"
	"errors"
	"fmt"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"

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
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("%w: no mother profile associated with this account", repository.ErrForbidden)
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

// GetChild returns the simple child detail (Endpoint: 12 - Get lightweight child profile)
func (s *ChildService) GetChild(ctx context.Context, userID string, childID uuid.UUID) (*models.ChildSimpleResponse, error) {
	return s.repo.GetChild(ctx, childID)
}

// GetListChild returns the lightweight children list (Endpoint: 11)
func (s *ChildService) GetListChild(ctx context.Context, userID string) ([]models.ChildListItem, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID formatting")
	}

	children, err := s.repo.GetListChild(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	for i, c := range children {
		children[i].Age = utils.FormatAgeString(c.BirthDate)
	}

	return children, nil
}
