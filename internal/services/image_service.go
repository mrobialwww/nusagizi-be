package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"nusagizi_be/internal/infrastructure/storage"
	"nusagizi_be/internal/repository"
)

type ImageService struct {
	storage       storage.ObjectStorage
	motherRepo    *repository.MotherProfileRepository
	childRepo     *repository.ChildRepository
	caregiverRepo *repository.CaregiverRepository
}

func NewImageService(
	storage storage.ObjectStorage,
	motherRepo *repository.MotherProfileRepository,
	childRepo *repository.ChildRepository,
	caregiverRepo *repository.CaregiverRepository,
) *ImageService {
	return &ImageService{
		storage:       storage,
		motherRepo:    motherRepo,
		childRepo:     childRepo,
		caregiverRepo: caregiverRepo,
	}
}

func (s *ImageService) GeneratePresignPutURL(ctx context.Context, category, ownerID, userID string, childID *string, contentType string) (uploadURL, objectKey string, err error) {
	switch category {
	case "profile":
		ownerID = userID
	case "social", "child-profile":
		if childID != nil && *childID != "" {
			cID, parseErr := uuid.Parse(*childID)
			if parseErr != nil {
				return "", "", fmt.Errorf("invalid child_id format")
			}

			motherProfileID, fetchErr := s.childRepo.GetMotherProfileID(ctx, cID)
			if fetchErr != nil {
				return "", "", fmt.Errorf("child not found")
			}

			// Validate ownership: Is user the mother?
			err = checkMotherOwnership(ctx, userID, cID, s.motherRepo, s.childRepo)
			if err != nil {
				// Is user a caregiver with active engagement for this child?
				caregiverProfileID, cgErr := s.caregiverRepo.GetByUserID(ctx, userID)
				if cgErr != nil {
					return "", "", fmt.Errorf("forbidden: user has no access to this child")
				}

				hasAccess, checkErr := s.caregiverRepo.CheckCaregiverAccess(ctx, cID, caregiverProfileID)
				if checkErr != nil || !hasAccess {
					return "", "", fmt.Errorf("forbidden: user has no access to this child")
				}
			}

			// Valid access, assign ownerID to mother_profile_id
			ownerID = motherProfileID.String()
		} else if category == "child-profile" {
			// fallback: finding mother_profile_id from user_id if creating new child profile
			motherProfileID, fetchErr := s.motherRepo.GetByUserID(ctx, userID)
			if fetchErr != nil {
				return "", "", fmt.Errorf("forbidden: mother profile not found")
			}
			ownerID = motherProfileID.String()
		}
	}

	if ownerID == "" {
		ownerID = uuid.NewString()
	}

	if category == "" {
		category = "misc"
	}

	ext := "webp"
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		ext = "jpg"
	} else if strings.Contains(contentType, "png") {
		ext = "png"
	}

	objectKey = fmt.Sprintf("%s/%s/%s.%s", category, ownerID, uuid.NewString(), ext)

	uploadURL, err = s.storage.PresignPutURL(ctx, objectKey, contentType, 10*time.Minute)
	return uploadURL, objectKey, err
}
