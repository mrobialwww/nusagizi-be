package services

import (
	"context"
	"fmt"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

// checkChildAccess verifies if the given user owns the child (as mother) or has an active engagement (as caregiver).
func checkChildAccess(
	ctx context.Context,
	userID string,
	childID uuid.UUID,
	motherRepo *repository.MotherProfileRepository,
	caregiverRepo *repository.CaregiverRepository,
	childRepo *repository.ChildRepository,
) error {
	// Try mother profile
	motherProfileID, err := motherRepo.GetByUserID(ctx, userID)
	if err == nil {
		owns, err := childRepo.CheckOwnership(ctx, childID, motherProfileID)
		if err == nil && owns {
			return nil
		}
	}

	// Try caregiver profile
	caregiverProfileID, err := caregiverRepo.GetByUserID(ctx, userID)
	if err == nil {
		access, err := caregiverRepo.CheckCaregiverAccess(ctx, childID, caregiverProfileID)
		if err == nil && access {
			return nil
		}
	}

	return fmt.Errorf("%w: user does not have access to this child", repository.ErrForbidden)
}

// checkMotherOwnership verifies if the given user is the mother (owner) of the child.
// This is used for strict operations like UpdateChild or DeleteChild.
func checkMotherOwnership(
	ctx context.Context,
	userID string,
	childID uuid.UUID,
	motherRepo *repository.MotherProfileRepository,
	childRepo *repository.ChildRepository,
) error {
	motherProfileID, err := motherRepo.GetByUserID(ctx, userID)
	if err == nil {
		owns, err := childRepo.CheckOwnership(ctx, childID, motherProfileID)
		if err == nil && owns {
			return nil
		}
	}
	return fmt.Errorf("%w: user does not have owner access to this child", repository.ErrForbidden)
}

// checkCaregiverAccess verifies if the given user is a caregiver who has active access to the child.
func checkCaregiverAccess(
	ctx context.Context,
	userID string,
	childID uuid.UUID,
	caregiverRepo *repository.CaregiverRepository,
) error {
	caregiverProfileID, err := caregiverRepo.GetByUserID(ctx, userID)
	if err == nil {
		hasAccess, err := caregiverRepo.CheckCaregiverAccess(ctx, childID, caregiverProfileID)
		if err == nil && hasAccess {
			return nil
		}
	}
	return fmt.Errorf("%w: user does not have caregiver access to this child", repository.ErrForbidden)
}
