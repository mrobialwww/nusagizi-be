package services

import (
	"context"
	"fmt"

	"nusagizi_be/internal/models/photos_contacts"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"

	"github.com/google/uuid"
)

type PhotosContactsService struct {
	repo          *repository.PhotosContactsRepository
	motherRepo    *repository.MotherProfileRepository
	childRepo     *repository.ChildRepository
	caregiverRepo *repository.CaregiverRepository
	r2PublicURL   string
}

func NewPhotosContactsService(
	repo *repository.PhotosContactsRepository,
	motherRepo *repository.MotherProfileRepository,
	childRepo *repository.ChildRepository,
	caregiverRepo *repository.CaregiverRepository,
	r2PublicURL string,
) *PhotosContactsService {
	return &PhotosContactsService{
		repo:          repo,
		motherRepo:    motherRepo,
		childRepo:     childRepo,
		caregiverRepo: caregiverRepo,
		r2PublicURL:   r2PublicURL,
	}
}

// GetContacts (Endpoint: 41)
func (s *PhotosContactsService) GetContacts(ctx context.Context, userID string) ([]photos_contacts.ContactResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	contacts, err := s.repo.GetContacts(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}
	for i := range contacts {
		contacts[i].PhotoURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, contacts[i].PhotoURL)
	}
	return contacts, nil
}

// AddContact (endpoint: 49)
func (s *PhotosContactsService) AddContact(ctx context.Context, userID string, relatedMotherProfileID uuid.UUID) (uuid.UUID, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}
	if motherProfileID == relatedMotherProfileID {
		return uuid.Nil, fmt.Errorf("cannot add yourself as contact")
	}
	return s.repo.CreateContact(ctx, motherProfileID, relatedMotherProfileID)
}

// DeleteContact (Endpoint: 42)
func (s *PhotosContactsService) DeleteContact(ctx context.Context, userID string, contactID uuid.UUID) error {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteContact(ctx, contactID, motherProfileID)
}

// GetMotherChildPhotos (Endpoint: 43)
func (s *PhotosContactsService) GetMotherChildPhotos(ctx context.Context, userID string, childID *uuid.UUID, latestPerChild bool) ([]photos_contacts.ChildPhotoResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	photos, err := s.repo.GetMotherChildPhotos(ctx, motherProfileID, childID, latestPerChild)
	if err != nil {
		return nil, err
	}
	for i := range photos {
		if resolved := utils.ResolvePublicPhotoURL(s.r2PublicURL, &photos[i].PhotoURL); resolved != nil {
			photos[i].PhotoURL = *resolved
		}
	}
	return photos, nil
}

// GetContactChildPhotos (Endpoint: 44)
func (s *PhotosContactsService) GetContactChildPhotos(ctx context.Context, userID string, contactID uuid.UUID) ([]photos_contacts.ChildPhotoResponse, error) {
	// The caller must own the contact.
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if contactID belongs to motherProfileID
	contacts, err := s.repo.GetContacts(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}
	owns := false
	for _, c := range contacts {
		if c.ContactID == contactID {
			owns = true
			break
		}
	}
	if !owns {
		return nil, fmt.Errorf("%w: contact does not belong to you", repository.ErrForbidden)
	}

	// Fetch photos of children from the contact
	photos, err := s.repo.GetContactChildPhotos(ctx, contactID)
	if err != nil {
		return nil, err
	}
	for i := range photos {
		if resolved := utils.ResolvePublicPhotoURL(s.r2PublicURL, &photos[i].PhotoURL); resolved != nil {
			photos[i].PhotoURL = *resolved
		}
	}
	return photos, nil
}

// GetAllChildPhotos (Endpoint: 45)
func (s *PhotosContactsService) GetAllChildPhotos(ctx context.Context, userID string) ([]photos_contacts.ChildPhotoResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	photos, err := s.repo.GetAllChildPhotos(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}
	for i := range photos {
		if resolved := utils.ResolvePublicPhotoURL(s.r2PublicURL, &photos[i].PhotoURL); resolved != nil {
			photos[i].PhotoURL = *resolved
		}
	}
	return photos, nil
}

// GetPhotoDetail (Endpoint: 46)
func (s *PhotosContactsService) GetPhotoDetail(ctx context.Context, userID string, photoID uuid.UUID) (*photos_contacts.ChildPhotoResponse, error) {
	// 1. Get the photo to know the childID
	photo, err := s.repo.GetPhotoDetail(ctx, photoID)
	if err != nil {
		return nil, err
	}

	// 2. Check if user is the mother
	if err := checkMotherOwnership(ctx, userID, photo.ChildID, s.motherRepo, s.childRepo); err == nil {
		return photo, nil
	}

	// 3. Check if user is an active caregiver
	if err := checkCaregiverAccess(ctx, userID, photo.ChildID, s.caregiverRepo); err == nil {
		return photo, nil
	}

	// 4. Check if user is an authorized contact
	hasAccess, err := s.repo.CheckContactPhotoAccess(ctx, userID, photoID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("%w: user does not have permission to view this photo", repository.ErrForbidden)
	}

	if resolved := utils.ResolvePublicPhotoURL(s.r2PublicURL, &photo.PhotoURL); resolved != nil {
		photo.PhotoURL = *resolved
	}
	return photo, nil
}

// AddPhotoMother (Endpoint: 47)
func (s *PhotosContactsService) AddPhotoMother(ctx context.Context, userID string, childID uuid.UUID, input *photos_contacts.CreatePhotoInput) (uuid.UUID, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}

	if input.Caption == "" {
		return uuid.Nil, fmt.Errorf("caption is required")
	}

	if input.Visibility == "selected_only" && len(input.ListVisibility) == 0 {
		return uuid.Nil, fmt.Errorf("list_visibility is required when visibility is 'selected_only'")
	}

	// Force is_review_required to false
	input.IsReviewRequired = false

	return s.repo.CreatePhotoWithStreak(ctx, childID, input)
}

// AddPhotoCaregiver (Endpoint: 48)
func (s *PhotosContactsService) AddPhotoCaregiver(ctx context.Context, userID string, childID uuid.UUID, input *photos_contacts.CreatePhotoInput) (uuid.UUID, error) {
	if err := checkCaregiverAccess(ctx, userID, childID, s.caregiverRepo); err != nil {
		return uuid.Nil, err
	}

	if !input.IsReviewRequired {
		return uuid.Nil, fmt.Errorf("is_review_required must be true for caregiver uploads")
	}

	input.Visibility = "private" // Default for caregiver

	return s.repo.CreatePhoto(ctx, childID, input)
}

// UpdatePhoto (endpoint: 50)
func (s *PhotosContactsService) UpdatePhoto(ctx context.Context, userID string, childID uuid.UUID, photoID uuid.UUID, input *photos_contacts.UpdatePhotoInput) error {
	// Access control: photo must belong to mother's child
	// Fetch photo to get childID
	photo, err := s.repo.GetPhotoDetail(ctx, photoID)
	if err != nil {
		return err
	}
	if photo.ChildID != childID {
		return fmt.Errorf("%w: photo does not belong to the specified child", repository.ErrNotFound)
	}
	if err := checkMotherOwnership(ctx, userID, photo.ChildID, s.motherRepo, s.childRepo); err != nil {
		return err
	}

	if input.Visibility != nil && *input.Visibility == "selected_only" {
		if input.ListVisibility == nil || len(*input.ListVisibility) == 0 {
			return fmt.Errorf("list_visibility is required when visibility is 'selected_only'")
		}
	}

	return s.repo.UpdatePhoto(ctx, photoID, input)
}

// DeletePhoto (endpoint: 51)
func (s *PhotosContactsService) DeletePhoto(ctx context.Context, userID string, photoID uuid.UUID) error {
	photo, err := s.repo.GetPhotoDetail(ctx, photoID)
	if err != nil {
		return err
	}
	if err := checkMotherOwnership(ctx, userID, photo.ChildID, s.motherRepo, s.childRepo); err != nil {
		return err
	}

	return s.repo.DeletePhoto(ctx, photoID)
}
