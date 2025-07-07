package link

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FromForm populates a new Link from HTTP POST form input values and performs
// validation on all fields.
//
// It supports two expiration strategies:
//   - Absolute expiration time via `expires-at` (RFC3339 string)
//   - Relative expiration via `expires-in` (number of days from `now`)
//
// If both `expires-at` and `expires-in` are provided, `expires-at` takes
// precedence.
//
// Returns an error if required fields are missing or invalid.
func (l *Link) FromForm(form url.Values, now time.Time) (*Link, error) {
	// Required fields
	target := form.Get("target")
	slugLenStr := form.Get("slug-len")
	slugCharset := form.Get("slug-charset")
	expiresAtStr := form.Get("expires-at")
	expiresInStr := form.Get("expires-in") // In number of days

	// Optional fields
	password := form.Get("password")
	maxHitsStr := form.Get("max-hits")
	validFromStr := form.Get("valid-from")

	var err error

	err = validateTarget(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target: %w", err)
	}

	var passwordHash *string
	if password != "" {
		hashed, err := hashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %w", err)
		}
		passwordHash = &hashed
	}

	var expiresAt time.Time
	if expiresAtStr != "" {
		expiresAt, err = time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-at: %w", err)
		}
	} else if expiresInStr != "" {
		expiresIn, err := strconv.Atoi(expiresInStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expires-in: %w", err)
		}
		expiresAt = now.Add(time.Duration(expiresIn) * 24 * time.Hour)
	} else {
		return nil, errors.New("expires-at and expires-in are both missing")
	}

	err = validateExpiresAt(expiresAt, now)
	if err != nil {
		return nil, fmt.Errorf("invalid expires-at: %w", err)
	}

	var validFrom *time.Time
	if validFromStr != "" {
		t, err := time.Parse(time.RFC3339, validFromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid valid-from: %w", err)
		}
		validFrom = &t
	}

	err = validateValidFrom(validFrom, now)
	if err != nil {
		return nil, fmt.Errorf("invalid valid-from: %w", err)
	}

	adminExpiresAt := expiresAt.Add(24 * time.Hour)

	err = validateTimes(validFrom, expiresAt, adminExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("one or more invalid times: %w", err)
	}

	var maxHits *int
	if maxHitsStr != "" {
		i, err := strconv.Atoi(maxHitsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid max-hits: %w", err)
		}
		maxHits = &i
	}

	err = validateMaxHits(maxHits)
	if err != nil {
		return nil, fmt.Errorf("invalid max-hits: %w", err)
	}

	slugLen, err := strconv.Atoi(slugLenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid slug-len: %w", err)
	}

	slug, err := generateSlug(slugLen, strings.ToLower(slugCharset))
	if err != nil {
		return nil, fmt.Errorf("error generating slug: %w", err)
	}

	adminToken, err := generateAdminToken(adminTokenLen)
	if err != nil {
		return nil, fmt.Errorf("error generating admin token: %w", err)
	}

	link := &Link{
		Slug:           slug,
		AdminToken:     adminToken,
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

	return link, nil
}
