package services

import (
	"context"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo          *repository.UserRepository
	motherRepo    *repository.MotherProfileRepository
	caregiverRepo *repository.CaregiverRepository
}

func NewUserService(repo *repository.UserRepository, motherRepo *repository.MotherProfileRepository, caregiverRepo *repository.CaregiverRepository) *UserService {
	return &UserService{repo: repo, motherRepo: motherRepo, caregiverRepo: caregiverRepo}
}

// GetUser returns a user's profile (Endpoint: 3)
func (s *UserService) GetUser(ctx context.Context, userID string) (*models.UserProfileResponse, error) {
	return s.repo.GetProfileByID(ctx, userID)
}

// UpdateUser partially updates a user's profile (Endpoint: 1)
func (s *UserService) UpdateUser(ctx context.Context, userID string, input models.UpdateUserInput) error {
	return s.repo.Update(ctx, userID, input)
}

// SoftDeleteUser marks a user as deleted (Endpoint: 2)
func (s *UserService) SoftDeleteUser(ctx context.Context, userID string) error {
	return s.repo.SoftDelete(ctx, userID)
}

// CreateMotherProfile creates a mother profile for the user (Endpoint: 4)
func (s *UserService) CreateMotherProfile(ctx context.Context, userID string) (uuid.UUID, error) {
	return s.motherRepo.Create(ctx, userID)
}

// CreateCaregiverProfile creates a caregiver profile for the user (Endpoint: 5)
func (s *UserService) CreateCaregiverProfile(ctx context.Context, userID string) (uuid.UUID, error) {
	return s.caregiverRepo.Create(ctx, userID)
}

// GetMotherProfile returns the mother profile (Endpoint: 6)
func (s *UserService) GetMotherProfile(ctx context.Context, userID string) (*models.MotherProfileResponse, error) {
	// 1. Validate ownership & get ID
	actualMotherID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch user profile for join fields
	userProfile, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Assemble response
	return &models.MotherProfileResponse{
		ID:          actualMotherID.String(),
		UserID:      userProfile.ID,
		FullName:    userProfile.FullName,
		Email:       userProfile.Email,
		PhoneNumber: userProfile.PhoneNumber,
		PhotoURL:    userProfile.PhotoURL,
	}, nil
}

// GetCaregiverProfile returns the caregiver profile (endpoint: 55)
func (s *UserService) GetCaregiverProfile(ctx context.Context, userID string) (*models.CaregiverProfileResponse, error) {
	// 1. Validate ownership & get ID
	actualCaregiverID, err := s.caregiverRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch user profile for join fields
	userProfile, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Assemble response
	return &models.CaregiverProfileResponse{
		ID:          actualCaregiverID.String(),
		UserID:      userProfile.ID,
		FullName:    userProfile.FullName,
		Email:       userProfile.Email,
		PhoneNumber: userProfile.PhoneNumber,
		PhotoURL:    userProfile.PhotoURL,
	}, nil
}
