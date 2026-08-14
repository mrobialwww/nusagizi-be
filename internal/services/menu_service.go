package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"nusagizi_be/internal/config"
	"nusagizi_be/internal/models/menu"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// MenuService orchestrates menu generation: fetching child data, calling the AI engine, and persisting the resulting menu.
type MenuService struct {
	cfg            *config.Config
	menuRepo       *repository.MenuRepository
	motherRepo     *repository.MotherProfileRepository
	childRepo      *repository.ChildRepository
	childNutriRepo *repository.ChildNutritionRepository
	httpClient     *http.Client
}

func NewMenuService(cfg *config.Config, menuRepo *repository.MenuRepository, motherRepo *repository.MotherProfileRepository, childRepo *repository.ChildRepository, childNutriRepo *repository.ChildNutritionRepository) *MenuService {
	return &MenuService{
		cfg:            cfg,
		menuRepo:       menuRepo,
		motherRepo:     motherRepo,
		childRepo:      childRepo,
		childNutriRepo: childNutriRepo,
		httpClient:     &http.Client{Timeout: 30 * time.Second}, // Reuse client for HTTP connection pooling
	}
}

// GenerateMenu (Endpoint: 39) generates and saves today's menu for all children of a mother.
func (s *MenuService) GenerateMenu(ctx context.Context, userID string, reportDate time.Time) error {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("mother profile not found: %w", err)
		}
		return fmt.Errorf("failed to get mother profile: %w", err)
	}

	// Extract batch data via service method
	children, err := s.getChildrenDataForMenu(ctx, motherProfileID)
	if err != nil {
		return fmt.Errorf("failed to build AI payload: %w", err)
	}
	if len(children) == 0 {
		return errors.New("no children found for this mother")
	}

	// Encode Go structs to JSON bytes for HTTP body
	payloadBytes, err := json.Marshal(children)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Call AI engine with retries for transient failures
	aiResponse, err := s.callAIEngine(ctx, payloadBytes)
	if err != nil {
		return err
	}
	if len(aiResponse.Anak) == 0 {
		return errors.New("ai engine returned no menu data")
	}
	if len(aiResponse.Anak) != len(children) {
		slog.WarnContext(ctx, "ai engine returned a different number of children than requested",
			"mother_profile_id", motherProfileID,
			"requested", len(children),
			"returned", len(aiResponse.Anak),
		)
	}

	return s.saveAllChildMenus(ctx, aiResponse.Anak, reportDate)
}

// saveAllChildMenus saves children's menus concurrently (max 4).
// Errors from individual saves are collected and do not block others.
func (s *MenuService) saveAllChildMenus(ctx context.Context, childList []menu.FoodEngineAnak, today time.Time) error {
	sem := make(chan struct{}, 4) // maxConcurrentChildSaves
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, child := range childList {
		wg.Add(1)
		sem <- struct{}{}
		go func(child menu.FoodEngineAnak) {
			defer wg.Done()
			defer func() { <-sem }()

			childID, err := uuid.Parse(child.ID)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("child %s: invalid id: %w", child.ID, err))
				mu.Unlock()
				return
			}

			// Check if any recipe in today's report for this child is reused from a previous report
			existingCNR, err := s.childNutriRepo.GetReportByChildAndDate(ctx, childID, today)
			if err == nil && existingCNR != nil {

				isReused, err := s.childNutriRepo.IsReportReused(ctx, existingCNR.ID)
				if err == nil && isReused {
					// Skip this child because their menu is being reused, allow others to continue
					slog.WarnContext(ctx, "Skipping menu replacement for child due to reused report", "child_id", childID)
					return
				}
			}

			// Calculate Macros, Sum target macros across all sessions (breakfast, lunch, dinner) and snacks
			var calories, protein, fat, carbo float64
			for _, sesi := range child.Menu.Sesi {
				calories += sesi.Makro["energi"].Target
				protein += sesi.Makro["protein"].Target
				fat += sesi.Makro["lemak"].Target
				carbo += sesi.Makro["karbo"].Target
			}

			for _, selingan := range child.Menu.Selingan {
				calories += selingan.Makro.Energi
				protein += selingan.Makro.Protein
				fat += selingan.Makro.Lemak
				carbo += selingan.Makro.Karbo
			}

			// Save to Database
			if err := s.menuRepo.ReplaceTodayMenu(
				ctx, childID, today,
				calories, protein, fat, carbo,
				child.Menu,
			); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("child %s: failed to save new menu: %w", child.ID, err))
				mu.Unlock()
			}
		}(child)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("failed to save menu for %d child(ren): %w", len(errs), errors.Join(errs...))
	}
	return nil
}

// callAIEngine sends payload to AI engine and returns the decoded response.
func (s *MenuService) callAIEngine(ctx context.Context, payloadBytes []byte) (*menu.FoodEngineResponse, error) {
	url := s.cfg.AIServiceBaseURL + "/api/v1/menu/generate"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ai engine: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read ai engine response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ai engine returned status %d: %s", resp.StatusCode, string(body))
	}

	var aiResponse menu.FoodEngineResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		return nil, fmt.Errorf("ai engine returned an invalid response: %w", err)
	}
	return &aiResponse, nil
}

// getChildrenDataForMenu constructs the AI payload for all children of a mother.
// It uses batched queries for allergies, growth, and medical notes
func (s *MenuService) getChildrenDataForMenu(ctx context.Context, motherProfileID uuid.UUID) ([]menu.FoodEnginePayloadChild, error) {
	// Get all children data from mother
	children, err := s.childRepo.FetchChildren(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		return nil, nil
	}

	// Map children ID to child struct
	childIDs := make([]uuid.UUID, len(children))
	for i, c := range children {
		childIDs[i] = c.ID
	}

	// Get child allergies data in batch
	allergiesByChild, err := s.childRepo.FetchAllergies(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Get child growth data in batch
	growthByChild, err := s.childRepo.FetchGrowthHistory(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Get child medical notes data in batch
	medicalNotesByChild, err := s.childRepo.FetchActiveMedicalNotes(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Map medical notes ID to nutrition targets
	noteIDs := make([]uuid.UUID, 0, len(medicalNotesByChild))
	for _, note := range medicalNotesByChild {
		noteIDs = append(noteIDs, note.ID)
	}

	// Get nutrition targets data in batch
	targetsByNote, err := s.childRepo.FetchNutritionTargets(ctx, noteIDs)
	if err != nil {
		return nil, err
	}

	// Construct payload for each child
	payload := make([]menu.FoodEnginePayloadChild, 0, len(children))
	for _, c := range children {
		// Calculate age in months
		usiaBulan := utils.CalculateAgeInMonths(c.BirthDate, time.Now())

		// Map sex from DB format (male/female) to AI Engine format (boys/girls)
		sex := strings.ToLower(c.Sex)
		if sex == "male" || sex == "boy" || sex == "boys" || sex == "l" {
			sex = "boys"
		} else {
			sex = "girls"
		}

		// Create child payload (using c.ID.String() for Nama so AI Engine returns UUID in anak.ID)
		childPayload := menu.FoodEnginePayloadChild{
			ID:                   c.ID.String(),
			Nama:                 c.ID.String(),
			Sex:                  sex,
			UsiaBulan:            usiaBulan,
			PreferensiHarga:      "seimbang", // temporary hardcode
			JumlahProteinPerHari: 1,          // temporary hardcode
			Alergi:               allergiesByChild[c.ID],
			RiwayatBerat:         [][]any{},
			Kondisi:              []menu.FoodEnginePayloadKondisi{}, // temporary fill with empty array
		}

		// AI expects an empty array [] if there are no allergies
		if childPayload.Alergi == nil {
			childPayload.Alergi = []string{}
		}

		// Attach latest physical measurements and format weight history graph as [AgeInMonths, WeightKg]
		if growth, ok := growthByChild[c.ID]; ok && len(growth) > 0 {
			// Latest measurement is first in the slice (newest)
			latest := growth[0]
			childPayload.BeratKg = latest.WeightKg
			childPayload.PanjangCm = latest.HeightCm
			childPayload.LingkarKepalaCm = latest.HeadCircCm

			// Convert each measurement into a [AgeInMonths, WeightKg] point for the growth curve
			history := make([][]any, len(growth))
			for i, g := range growth {
				history[i] = []any{utils.CalculateAgeInMonths(c.BirthDate, g.MeasurementDate), g.WeightKg}
			}
			childPayload.RiwayatBerat = history
		}

		// Initialize default medical prescription (calculates max safe daily sodium based on child's age)
		resepDokter := menu.FoodEnginePayloadResepDokter{
			Tingkat:         nil, // temporary hardcode
			NomorStr:        nil, // temporary hardcode
			NamaDokter:      "-",
			CatatanTambahan: "-",
			KomposisiPerKg: menu.FoodEnginePayloadKomposisi{
				EnergiKkal:    0,
				ProteinG:      0,
				NatriumMgMaks: getNatriumMaks(usiaBulan),
			},
		}
		// Override default prescription if child has an active doctor's note and targeted nutrition
		if note, ok := medicalNotesByChild[c.ID]; ok {
			resepDokter.NamaDokter = note.DoctorName
			resepDokter.CatatanTambahan = note.Recommendation
			if targets, ok := targetsByNote[note.ID]; ok {
				resepDokter.KomposisiPerKg.EnergiKkal = targets.EnergiKkal
				resepDokter.KomposisiPerKg.ProteinG = targets.ProteinG
			}
		}
		childPayload.ResepDokter = resepDokter

		payload = append(payload, childPayload)
	}

	return payload, nil
}

func getNatriumMaks(usiaBulan int) float64 {
	if usiaBulan <= 5 {
		return 120
	} else if usiaBulan <= 11 {
		return 370
	} else if usiaBulan <= 36 {
		return 800
	}
	return 900
}
