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
	"sync"
	"time"

	"nusagizi_be/internal/config"
	"nusagizi_be/internal/models/menu"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// MenuService orchestrates menu generation: fetching child data, calling the AI engine, and persisting the resulting menu.
type MenuService struct {
	cfg        *config.Config
	menuRepo   *repository.MenuRepository
	motherRepo *repository.MotherProfileRepository
	childRepo  *repository.ChildRepository
	httpClient *http.Client
}

func NewMenuService(cfg *config.Config, menuRepo *repository.MenuRepository, motherRepo *repository.MotherProfileRepository, childRepo *repository.ChildRepository) *MenuService {
	return &MenuService{
		cfg:        cfg,
		menuRepo:   menuRepo,
		motherRepo: motherRepo,
		childRepo:  childRepo,
		httpClient: &http.Client{Timeout: 30 * time.Second}, // Reuse client for HTTP connection pooling
	}
}

// GenerateMenu generates and saves today's menu for all children of a mother.
func (s *MenuService) GenerateMenu(ctx context.Context, userID string) error {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("mother profile not found: %w", err)
		}
		return fmt.Errorf("failed to get mother profile: %w", err)
	}

	// Extract batch data directly from child repository (Service A -> Repo B pattern).
	children, err := s.childRepo.GetChildrenDataForMenu(ctx, motherProfileID)
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

	return s.saveAllChildMenus(ctx, aiResponse.Anak, time.Now())
}

// saveAllChildMenus saves children's menus concurrently (max 4).
// Errors from individual saves are collected and do not block others.
func (s *MenuService) saveAllChildMenus(ctx context.Context, anakList []menu.FoodEngineAnak, today time.Time) error {
	sem := make(chan struct{}, 4) // maxConcurrentChildSaves
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, anak := range anakList {
		wg.Add(1)
		sem <- struct{}{}
		go func(anak menu.FoodEngineAnak) {
			defer wg.Done()
			defer func() { <-sem }()

			childID, err := uuid.Parse(anak.ID)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("child %s: invalid id: %w", anak.ID, err))
				mu.Unlock()
				return
			}

			// 1. Calculate Macros, Sum target macros across all sessions (breakfast, lunch, dinner) and snacks
			var calories, protein, fat, carbo float64
			for _, sesi := range anak.Menu.Sesi {
				calories += sesi.Makro["energi"].Target
				protein += sesi.Makro["protein"].Target
				fat += sesi.Makro["lemak"].Target
				carbo += sesi.Makro["karbo"].Target
			}

			for _, selingan := range anak.Menu.Selingan {
				calories += selingan.Makro.Energi
				protein += selingan.Makro.Protein
				fat += selingan.Makro.Lemak
				carbo += selingan.Makro.Karbo
			}

			// 2. Save to Database
			if err := s.menuRepo.ReplaceTodayMenu(
				ctx, childID, today,
				calories, protein, fat, carbo,
				anak.Menu,
			); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("child %s: failed to save new menu: %w", anak.ID, err))
				mu.Unlock()
			}
		}(anak)
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

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ai engine returned status %d: %s", resp.StatusCode, string(body))
	}

	var aiResponse menu.FoodEngineResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		return nil, fmt.Errorf("ai engine returned an invalid response: %w", err)
	}
	return &aiResponse, nil
}
