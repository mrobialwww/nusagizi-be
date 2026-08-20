package repository

import (
	"context"
	"fmt"
	"time"

	"nusagizi_be/internal/models/photos_contacts"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotosContactsRepository struct {
	pool *pgxpool.Pool
}

func NewPhotosContactsRepository(pool *pgxpool.Pool) *PhotosContactsRepository {
	return &PhotosContactsRepository{pool: pool}
}

// GetContacts returns a list of contacts (friends) for a mother profile.
func (r *PhotosContactsRepository) GetContacts(ctx context.Context, motherProfileID uuid.UUID) ([]photos_contacts.ContactResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT 
			c.id AS contact_id, 
			c.related_mother_profile_id, 
			u.full_name, 
			u.photo_url
		FROM contacts c
		JOIN mother_profiles mp ON c.related_mother_profile_id = mp.id
		JOIN users u ON mp.user_id = u.id
		WHERE c.mother_profile_id = $1
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []photos_contacts.ContactResponse
	for rows.Next() {
		var contact photos_contacts.ContactResponse
		if err := rows.Scan(
			&contact.ContactID,
			&contact.RelatedMotherProfileID,
			&contact.FullName,
			&contact.PhotoURL,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, rows.Err()
}

// CreateContact creates a mutual contact relationship.
func (r *PhotosContactsRepository) CreateContact(ctx context.Context, motherProfileID, relatedMotherProfileID uuid.UUID) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// Check if already exists
	var exists bool
	queryCheck := `
		SELECT EXISTS(
			SELECT 1 
			FROM contacts 
			WHERE mother_profile_id = $1 
				AND related_mother_profile_id = $2
		)
	`
	err = tx.QueryRow(ctx, queryCheck, motherProfileID, relatedMotherProfileID).Scan(&exists)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, fmt.Errorf("contact already exists")
	}

	// Insert Contact
	queryInsert := `
		INSERT INTO contacts (id, mother_profile_id, related_mother_profile_id)
		VALUES ($1, $2, $3)
	`

	newID1 := uuid.New()
	_, err = tx.Exec(ctx, queryInsert, newID1, motherProfileID, relatedMotherProfileID)
	if err != nil {
		return uuid.Nil, err
	}

	newID2 := uuid.New()
	_, err = tx.Exec(ctx, queryInsert, newID2, relatedMotherProfileID, motherProfileID)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return newID1, nil
}

// DeleteContact deletes a contact relationship.
func (r *PhotosContactsRepository) DeleteContact(ctx context.Context, contactID uuid.UUID, requesterMotherProfileID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// We need to find the pair first, We'll delete both sides of the friendship.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Find mother profile id and related mother profile id
	var motherProfileID, relatedMotherProfileID uuid.UUID
	queryFind := `
		SELECT mother_profile_id, related_mother_profile_id 
		FROM contacts 
		WHERE id = $1 
			AND mother_profile_id = $2
	`
	err = tx.QueryRow(ctx, queryFind, contactID, requesterMotherProfileID).Scan(&motherProfileID, &relatedMotherProfileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("record not found")
		}
		return err
	}

	// Delete both directions
	queryDelete := `
		DELETE FROM contacts 
		WHERE (mother_profile_id, related_mother_profile_id) IN (
			($1, $2),
			($2, $1)
		)
	`
	_, err = tx.Exec(ctx, queryDelete, motherProfileID, relatedMotherProfileID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetMotherChildPhotos gets all photos for a mother's children.
func (r *PhotosContactsRepository) GetMotherChildPhotos(ctx context.Context, motherProfileID uuid.UUID, childID *uuid.UUID, latestPerChild bool) ([]photos_contacts.ChildPhotoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	motherChildPhotosQuery := `
		SELECT p.id, p.child_id, p.photo_url, p.is_review_required, p.created_at
		FROM child_photos p
		JOIN children c ON p.child_id = c.id
		WHERE c.mother_profile_id = $1
		ORDER BY p.created_at DESC
	`

	// motherChildDailyPhotoQuery returns one photo per day (the latest that
	// day) for a single child — used for daily/timeline views.
	motherChildDailyPhotoQuery := `
		SELECT DISTINCT ON (DATE(p.created_at AT TIME ZONE 'UTC'))
			p.id, p.child_id, p.photo_url, p.is_review_required, p.created_at
		FROM child_photos p
		JOIN children c ON p.child_id = c.id
		WHERE c.mother_profile_id = $1 AND c.id = $2
		ORDER BY DATE(p.created_at AT TIME ZONE 'UTC') DESC, p.created_at DESC
	`

	// motherChildLatestPerChildQuery returns the single latest photo per child across all children
	motherChildLatestPerChildQuery := `
		SELECT DISTINCT ON (c.id)
			p.id, p.child_id, p.photo_url, p.is_review_required, p.created_at
		FROM child_photos p
		JOIN children c ON p.child_id = c.id
		WHERE c.mother_profile_id = $1
		ORDER BY c.id, p.created_at DESC
	`

	var query string
	var args []any

	if childID != nil {
		query = motherChildDailyPhotoQuery
		args = []any{motherProfileID, *childID}
	} else if latestPerChild {
		query = motherChildLatestPerChildQuery
		args = []any{motherProfileID}
	} else {
		query = motherChildPhotosQuery
		args = []any{motherProfileID}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []photos_contacts.ChildPhotoResponse
	for rows.Next() {
		var p photos_contacts.ChildPhotoResponse
		if err := rows.Scan(&p.ID, &p.ChildID, &p.PhotoURL, &p.IsReviewRequired, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

// GetContactChildPhotos gets photos shared with a specific contact.
func (r *PhotosContactsRepository) GetContactChildPhotos(ctx context.Context, contactID uuid.UUID) ([]photos_contacts.ChildPhotoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT p.id, p.child_id, p.photo_url, p.created_at
		FROM child_photos p
		LEFT JOIN photo_shares psw ON psw.child_photo_id = p.id
		JOIN children c ON p.child_id = c.id
		JOIN contacts ct ON ct.id = $1
		WHERE 
			c.mother_profile_id = ct.related_mother_profile_id
			AND p.is_review_required = false
			AND (
				p.visibility = 'all' 
				OR (p.visibility = 'selected_only' 
				AND psw.contact_id = $1)
			)
		ORDER BY p.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []photos_contacts.ChildPhotoResponse
	for rows.Next() {
		var p photos_contacts.ChildPhotoResponse
		if err := rows.Scan(&p.ID, &p.ChildID, &p.PhotoURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

// GetAllChildPhotos gets all photos (mother's children + contacts' children).
func (r *PhotosContactsRepository) GetAllChildPhotos(ctx context.Context, motherProfileID uuid.UUID) ([]photos_contacts.ChildPhotoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// A UNION query could be used here, combining GetMotherChildPhotos and GetContactChildPhotos logic
	// but filtered for is_review_required = false.
	query := `
		SELECT DISTINCT p.id, p.child_id, p.photo_url, p.created_at
		FROM child_photos p
		JOIN children c ON p.child_id = c.id
		LEFT JOIN contacts ct ON ct.related_mother_profile_id = c.mother_profile_id 
			AND ct.mother_profile_id = $1
		LEFT JOIN photo_shares psw ON psw.child_photo_id = p.id
		WHERE 
			p.is_review_required = false
			AND (
				c.mother_profile_id = $1 
				OR (
					ct.id IS NOT NULL 
					AND (
						p.visibility = 'all' 
						OR (p.visibility = 'selected_only' AND psw.contact_id = ct.id)
					)
				)
			)
		ORDER BY p.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []photos_contacts.ChildPhotoResponse
	for rows.Next() {
		var p photos_contacts.ChildPhotoResponse
		if err := rows.Scan(&p.ID, &p.ChildID, &p.PhotoURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

// GetPhotoDetail gets photo details by ID.
func (r *PhotosContactsRepository) GetPhotoDetail(ctx context.Context, photoID uuid.UUID) (*photos_contacts.ChildPhotoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var p photos_contacts.ChildPhotoResponse
	query := `
		SELECT id, child_id, photo_url, caption, created_at
		FROM child_photos
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, photoID).Scan(&p.ID, &p.ChildID, &p.PhotoURL, &p.Caption, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	return &p, nil
}

// CheckContactPhotoAccess verifies if the given userID has contact-level access to the photo.
func (r *PhotosContactsRepository) CheckContactPhotoAccess(ctx context.Context, userID string, photoID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM child_photos p
			JOIN children c ON p.child_id = c.id
			JOIN mother_profiles mp_owner ON c.mother_profile_id = mp_owner.id
			JOIN contacts ct ON ct.mother_profile_id = mp_owner.id
			JOIN mother_profiles mp_viewer ON ct.related_mother_profile_id = mp_viewer.id
			LEFT JOIN photo_shares psw ON psw.child_photo_id = p.id
			WHERE p.id = $1
				AND mp_viewer.user_id = $2
				AND p.is_review_required = false
				AND (
					p.visibility = 'all'
					OR (p.visibility = 'selected_only' AND psw.contact_id = ct.id)
				)
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, query, photoID, userID).Scan(&exists)
	return exists, err
}

// CreatePhoto adds a new photo and its visibility list.
func (r *PhotosContactsRepository) CreatePhoto(ctx context.Context, childID uuid.UUID, input *photos_contacts.CreatePhotoInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	newID := uuid.New()

	// Insert photo
	queryInsert := `
		INSERT INTO child_photos (id, child_id, photo_url, caption, visibility, is_review_required)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, queryInsert, newID, childID, input.URL, input.Caption, input.Visibility, input.IsReviewRequired)
	if err != nil {
		return uuid.Nil, err
	}

	// Insert photo shared with
	if input.Visibility == "selected_only" && len(input.ListVisibility) > 0 {
		queryShared := `
			INSERT INTO photo_shares (child_photo_id, contact_id)
			VALUES ($1, $2)
		`
		for _, contactID := range input.ListVisibility {
			_, err = tx.Exec(ctx, queryShared, newID, contactID)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

// CreatePhotoWithStreak adds a new photo and its visibility list, and updates the child's streak.
func (r *PhotosContactsRepository) CreatePhotoWithStreak(ctx context.Context, childID uuid.UUID, input *photos_contacts.CreatePhotoInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Lock the child row
	var currentStreak int
	var lastStreakDate *time.Time
	queryLock := `
		SELECT streak_days, last_streak_date 
		FROM children 
		WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, queryLock, childID).Scan(&currentStreak, &lastStreakDate)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Calculate new streak securely
	today := time.Now().UTC().Truncate(24 * time.Hour)
	newStreak := currentStreak

	if lastStreakDate == nil {
		newStreak = 1
	} else {
		lastDate := lastStreakDate.UTC().Truncate(24 * time.Hour)
		diff := today.Sub(lastDate).Hours() / 24

		if diff == 1 {
			newStreak = currentStreak + 1
		} else if diff > 1 {
			newStreak = 1
		}
		// if diff == 0, newStreak remains currentStreak
	}

	newID := uuid.New()

	// 3. Insert photo
	queryInsert := `
		INSERT INTO child_photos (id, child_id, photo_url, caption, visibility, is_review_required)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, queryInsert, newID, childID, input.URL, input.Caption, input.Visibility, input.IsReviewRequired)
	if err != nil {
		return uuid.Nil, err
	}

	// 4. Insert photo shared with
	if input.Visibility == "selected_only" && len(input.ListVisibility) > 0 {
		queryShared := `
			INSERT INTO photo_shares (child_photo_id, contact_id)
			VALUES ($1, $2)
		`
		for _, contactID := range input.ListVisibility {
			_, err = tx.Exec(ctx, queryShared, newID, contactID)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	// 5. Update children streak
	queryUpdateChild := `
		UPDATE children
		SET streak_days = $1, last_streak_date = $2
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, queryUpdateChild, newStreak, today, childID)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

// UpdatePhoto updates a photo's details using Sync/UPSERT for shared list.
func (r *PhotosContactsRepository) UpdatePhoto(ctx context.Context, photoID uuid.UUID, input *photos_contacts.UpdatePhotoInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update child_photos fields dynamically
	queryUpdate := `
		UPDATE child_photos
		SET
			caption = COALESCE($1, caption),
			visibility = COALESCE($2, visibility),
			is_review_required = COALESCE($3, is_review_required)
		WHERE id = $4
	`
	res, err := tx.Exec(ctx, queryUpdate, input.Caption, input.Visibility, input.IsReviewRequired, photoID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}

	// Sync visibility list
	if input.ListVisibility != nil {
		if len(*input.ListVisibility) == 0 {

			// Delete all existing shared photos
			queryDeleteAll := `
				DELETE FROM photo_shares 
				WHERE child_photo_id = $1`

			_, err = tx.Exec(ctx, queryDeleteAll, photoID)
			if err != nil {
				return err
			}
		} else {

			// Delete all existing shared photos
			queryDeleteUnchecked := `
				DELETE FROM photo_shares 
				WHERE child_photo_id = $1 
					AND contact_id != ALL($2::uuid[])
			`

			// pgx array encoding workaround for []uuid.UUID
			var listVisibilityStrings []string
			for _, id := range *input.ListVisibility {
				listVisibilityStrings = append(listVisibilityStrings, id.String())
			}

			_, err = tx.Exec(ctx, queryDeleteUnchecked, photoID, listVisibilityStrings)
			if err != nil {
				return err
			}

			// Insert shared photos
			queryInsertShared := `
				INSERT INTO photo_shares (child_photo_id, contact_id)
				VALUES ($1, $2)
				ON CONFLICT (child_photo_id, contact_id) DO NOTHING
			`
			for _, contactID := range *input.ListVisibility {
				_, err = tx.Exec(ctx, queryInsertShared, photoID, contactID)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit(ctx)
}

// DeletePhoto performs a hard delete.
func (r *PhotosContactsRepository) DeletePhoto(ctx context.Context, photoID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		DELETE FROM child_photos 
		WHERE id = $1`

	res, err := r.pool.Exec(ctx, query, photoID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}
