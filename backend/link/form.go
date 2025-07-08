package link

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrMissingRequiredFields = errors.New("missing one or more required fields: target, slug-length, slug-charset")
	ErrErrorHashingPassword  = errors.New("error hashing password")
	ErrAdminExpiresAtUpdate  = errors.New("admin_expires_at can only be updated if expires_at is also updated")
)

// Form field names expected in POST input
const (
	FormTarget      = "target"       // Required: the destination URL
	FormSlugLength  = "slug-length"  // Required: length of generated slug
	FormSlugCharset = "slug-charset" // Required: allowed characters for slug generation
	FormExpiresIn   = "expires-in"   // Optional: relative expiration (in days)
	FormExpiresAt   = "expires-at"   // Optional: absolute expiration time (RFC3339)
	FormPassword    = "password"     // Optional: password to protect link
	FormMaxHits     = "max-hits"     // Optional: max allowed hits (int)
	FormValidFrom   = "valid-from"   // Optional: RFC3339 start time for link validity
)

// FromForm creates and validates a new Link from POST form values.
//
// Fields:
//   - Required: target, slug-length, slug-charset
//   - Optional expiration: either `expires-at` (RFC3339) or `expires-in` (days)
//   - Optional: password, max-hits, valid-from (RFC3339)
//
// Returns a Validated link or an error if any fields are missing or invalid.
func FromForm(form url.Values, now time.Time) (*Validated, error) {
	// Required fields
	target := form.Get(FormTarget)
	slugLenStr := form.Get(FormSlugLength)
	slugCharset := form.Get(FormSlugCharset)
	expiresAtStr := form.Get(FormExpiresAt)
	expiresInStr := form.Get(FormExpiresIn)

	// Optional fields
	password := form.Get(FormPassword)
	maxHitsStr := form.Get(FormMaxHits)
	validFromStr := form.Get(FormValidFrom)

	if target == "" || slugLenStr == "" || slugCharset == "" {
		return nil, ErrMissingRequiredFields
	}

	var passwordHash *string
	if password != "" {
		hashed, err := hashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrErrorHashingPassword, err)
		}
		passwordHash = &hashed
	}

	var err error

	var expiresAt time.Time
	switch {
	case expiresAtStr != "":
		expiresAt, err = time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-at: %w", err)
		}
	case expiresInStr != "":
		expiresIn, err := strconv.Atoi(expiresInStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-in: %w", err)
		}
		expiresAt = now.Add(time.Duration(expiresIn) * 24 * time.Hour)
	default:
		return nil, errors.New("either expires-at or expires-in must be provided")
	}

	var validFrom *time.Time
	if validFromStr != "" {
		vf, err := time.Parse(time.RFC3339, validFromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid valid-from: %w", err)
		}
		validFrom = &vf
	}

	adminExpiresAt := expiresAt.Add(24 * time.Hour)

	var maxHits *int
	if maxHitsStr != "" {
		mh, err := strconv.Atoi(maxHitsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid max-hits: %w", err)
		}
		maxHits = &mh
	}

	link := &Link{
		Slug:           "",
		AdminToken:     "",
		Target:         target,
		PasswordHash:   passwordHash,
		MaxHits:        maxHits,
		ValidFrom:      validFrom,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      expiresAt,
		AdminExpiresAt: adminExpiresAt,
		HitCount:       0,
		SchemaVersion:  schemaVersion,
	}

	validated, err := Validate(link, now)
	if err != nil {
		return nil, fmt.Errorf("one or more form values are invalid: %w", err)
	}

	slugLen, err := strconv.Atoi(slugLenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid slug-length: %w", err)
	}

	slug, err := generateSlug(slugLen, strings.ToLower(slugCharset))
	if err != nil {
		return nil, fmt.Errorf("error generating slug: %w", err)
	}
	validated.SetSlug(slug)

	adminToken, err := generateAdminToken(adminTokenLen)
	if err != nil {
		return nil, fmt.Errorf("error generating admin token: %w", err)
	}
	validated.SetAdminToken(adminToken)

	return validated, nil
}

// PatchFromForm builds a PatchLink from form values for partial updates.
// It respects "remove-*" checkboxes by setting Field.Remove = true.
func PatchFromForm(form url.Values, original *Link) (*ValidatedPatch, error) {
	now := time.Now()

	patch := NewPatchLink(now)

	if expiresAtStr := form.Get(FormExpiresAt); expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-at format: %w", err)
		}
		patch.ExpiresAt = &expiresAt
	}

	if form.Has("remove-" + FormMaxHits) {
		patch.MaxHits.Remove = true
	} else if maxHitsStr := form.Get(FormMaxHits); maxHitsStr != "" {
		maxHits, err := strconv.Atoi(maxHitsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid max-hits format: %w", err)
		}
		patch.MaxHits.Value = &maxHits
	}

	if form.Has("remove-" + FormValidFrom) {
		patch.ValidFrom.Remove = true
	} else if validFromStr := form.Get(FormValidFrom); validFromStr != "" {
		validFrom, err := time.Parse(time.RFC3339, validFromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid valid-from format: %w", err)
		}
		patch.ValidFrom.Value = &validFrom
	}

	if form.Has("remove-" + FormPassword) {
		patch.PasswordHash.Remove = true
	} else if password := form.Get(FormPassword); password != "" {
		hash, err := hashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrErrorHashingPassword, err)
		}
		patch.PasswordHash.Value = &hash
	}

	validated, err := ValidatePatch(original, patch, now)
	if err != nil {
		return nil, err
	}

	return validated, nil
}
