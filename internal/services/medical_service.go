package services

import (
	"context"
	"fmt"

	"nusagizi_be/internal/models/medical"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

type MedicalService struct {
	repo       *repository.MedicalRepository
	motherRepo *repository.MotherProfileRepository
	childRepo  *repository.ChildRepository
}

func NewMedicalService(
	repo *repository.MedicalRepository,
	motherRepo *repository.MotherProfileRepository,
	childRepo *repository.ChildRepository,
) *MedicalService {
	return &MedicalService{
		repo:       repo,
		motherRepo: motherRepo,
		childRepo:  childRepo,
	}
}

// GetMedicalNotes (Endpoint: 60 & 61)
func (s *MedicalService) GetMedicalNotes(ctx context.Context, userID string, status string, month *int) ([]medical.MedicalNoteResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if status != "active" && status != "history" {
		return nil, fmt.Errorf("status must be active or history")
	}

	if month != nil && (*month < 1 || *month > 12) {
		return nil, fmt.Errorf("month must be between 1 and 12")
	}

	return s.repo.GetMedicalNotes(ctx, motherProfileID, status, month)
}

// GetMedicalNoteDetail (Endpoint: 62)
func (s *MedicalService) GetMedicalNoteDetail(ctx context.Context, userID string, noteID uuid.UUID) (*medical.MedicalNoteDetailResponse, error) {
	childID, err := s.repo.GetChildIDByNoteID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetMedicalNoteDetail(ctx, noteID)
}

// CreateMedicalNote (Endpoint: 63)
func (s *MedicalService) CreateMedicalNote(ctx context.Context, userID string, input *medical.CreateMedicalNoteInput) (uuid.UUID, error) {
	if err := checkMotherOwnership(ctx, userID, input.ChildID, s.motherRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}

	// Required fields validation
	if input.DoctorName == "" || input.Recommendation == "" || input.ValidDate == "" {
		return uuid.Nil, fmt.Errorf("doctor_name, recommendation, valid_date are required")
	}

	return s.repo.CreateMedicalNote(ctx, input)
}

// UpdateMedicalNote (Endpoint: 64)
func (s *MedicalService) UpdateMedicalNote(ctx context.Context, userID string, noteID uuid.UUID, input *medical.UpdateMedicalNoteInput) error {
	childID, err := s.repo.GetChildIDByNoteID(ctx, noteID)
	if err != nil {
		return err
	}
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateMedicalNote(ctx, noteID, input)
}

// DeleteMedicalNote (Endpoint: 65)
func (s *MedicalService) DeleteMedicalNote(ctx context.Context, userID string, noteID uuid.UUID) error {
	childID, err := s.repo.GetChildIDByNoteID(ctx, noteID)
	if err != nil {
		return err
	}
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.DeleteMedicalNote(ctx, noteID)
}
