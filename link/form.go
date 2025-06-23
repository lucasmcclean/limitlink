package link

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// schemaVersion represents the current schema version stored in each Link.
const schemaVersion int = 1

// Field names for form input keys.
const (
	fieldTarget      = "target"
	fieldExpiresIn   = "expires-in"
	fieldExpiresAt   = "expires-at"
	fieldValidFrom   = "valid-from"
	fieldMaxHits     = "max-hits"
	fieldSlugLength  = "slug-length"
	fieldSlugCharset = "slug-charset"
	fieldPassword    = "password"
)

// FromForm parses a Link object from URL-encoded form values.
// It performs all validation and sets defaults where appropriate.
func FromForm(form url.Values) (*Link, error) {
	now := time.Now()

	target, err := extractTarget(form)
	if err != nil {
		return nil, err
	}

	expiresAt, err := extractExpiresAtOrIn(form, now)
	if err != nil {
		return nil, err
	}

	validFrom, err := extractValidFrom(form, now, expiresAt)
	if err != nil {
		return nil, err
	}

	maxHits, err := extractMaxHits(form)
	if err != nil {
		return nil, err
	}

	slugLength, slugCharset, err := extractSlugConfig(form)
	if err != nil {
		return nil, err
	}

	slug, err := generateSlug(slugLength, slugCharset)
	if err != nil {
		return nil, errors.New("Error generating the slug, please try again.")
	}

	adminToken, err := generateAdminToken()
	if err != nil {
		return nil, errors.New("Error generating the admin token, please try again.")
	}

	passwordHash, err := extractPasswordHash(form)
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:            primitive.NewObjectID(),
		CreatedAt:     now,
		UpdatedAt:     now,
		HitCount:      0,
		Slug:          slug,
		Target:        target,
		ExpiresAt:     expiresAt,
		MaxHits:       maxHits,
		ValidFrom:     validFrom,
		AdminToken:    adminToken,
		PasswordHash:  passwordHash,
		SchemaVersion: schemaVersion,
	}, nil
}

// extractExpiresAtOrIn returns a time.Time based on either "expires_in" (days from now)
// or "expires_at" (an RFC3339 timestamp). Prefers "expires_in" if both are provided.
func extractExpiresAtOrIn(form url.Values, now time.Time) (time.Time, error) {
	value := getValue(fieldExpiresIn, form)
	if value != "" {
		days, err := strconv.Atoi(value)
		if err != nil || days <= 0 {
			return time.Time{}, errors.New("Expiration (in days) must be a whole number greater than 0.")
		}
		expiresAt := now.Add(time.Hour * 24 * time.Duration(days))
		if err := validateExpiresAt(expiresAt, now); err != nil {
			return time.Time{}, err
		}
		return expiresAt, nil
	}

	value = getValue(fieldExpiresAt, form)
	if value == "" {
		return time.Time{}, errors.New("Please specify when the link should expire either in days or with a date.")
	}
	expiresAt, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, errors.New("Expiration date must be in a valid format (e.g., 2025-12-31T23:59:59Z).")
	}
	if err := validateExpiresAt(expiresAt, now); err != nil {
		return time.Time{}, err
	}
	return expiresAt, nil
}

// extractValidFrom returns a pointer to a time.Time based on the "valid_from" field.
// Returns nil if not provided. Ensures the value is in the future and before expiry.
func extractValidFrom(form url.Values, now time.Time, expiresAt time.Time) (*time.Time, error) {
	value := getValue(fieldValidFrom, form)
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errors.New("Start date must be in a valid format (e.g., 2025-06-01T00:00:00Z).")
	}
	if err := validateValidFrom(&t, now); err != nil {
		return nil, err
	}
	if err := validateValidFromWithExpiresAt(&t, expiresAt); err != nil {
		return nil, err
	}
	return &t, nil
}

// extractMaxHits parses and validates the "max_hits" form value.
// Returns nil if not provided. Value must be greater than 0 if present.
func extractMaxHits(form url.Values) (*int, error) {
	value := getValue(fieldMaxHits, form)
	if value == "" {
		return nil, nil
	}
	mh, err := strconv.Atoi(value)
	if err != nil {
		return nil, errors.New("Maximum hits must be a valid whole number greater than 0.")
	}
	if mh <= 0 {
		return nil, errors.New("Maximum hits must be greater than 0.")
	}
	if err := validateMaxHits(&mh); err != nil {
		return nil, err
	}
	return &mh, nil
}

// extractTarget retrieves and validates the required "target" field (the destination URL).
func extractTarget(form url.Values) (string, error) {
	value := getValue(fieldTarget, form)
	if value == "" {
		return "", errors.New("You must enter a target URL.")
	}
	if err := validateTarget(value); err != nil {
		return "", err
	}
	return value, nil
}

// extractSlugConfig parses optional slug length and charset settings from the form.
// Applies defaults if not provided and enforces bounds on length.
func extractSlugConfig(form url.Values) (int, string, error) {
	slugLength := 7
	value := getValue(fieldSlugLength, form)
	if value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < minSlugLen || parsed > maxSlugLen {
			return 0, "", fmt.Errorf("Slug length must be a number between %d and %d.", minSlugLen, maxSlugLen)
		}
		slugLength = parsed
	}

	charset := getValue(fieldSlugCharset, form)
	switch charset {
	case "", "alphanumeric", "letters", "numbers":
	default:
		return 0, "", errors.New("Slug character set must be one of: 'alphanumeric', 'letters', or 'numbers'.")
	}

	return slugLength, charset, nil
}

// extractPasswordHash generates a hash for the "password" field, if provided.
func extractPasswordHash(form url.Values) (*string, error) {
	value := getValue(fieldPassword, form)
	if value == "" {
		return nil, nil
	}
	hash, err := generateHash(value)
	if err != nil {
		return nil, errors.New("Error generating the password hash, please try again.")
	}
	return &hash, nil
}

// getValue retrieves a form field and trims surrounding whitespace.
func getValue(fieldName string, form url.Values) string {
	return strings.TrimSpace(form.Get(fieldName))
}
