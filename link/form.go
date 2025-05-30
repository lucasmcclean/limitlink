package link

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	fieldTarget    = "target"
	fieldExpiresIn = "expires_in"
	fieldExpiresAt = "expires_at"
	fieldMaxHits   = "max_hits"
	fieldValidFrom = "valid_from"
	fieldPassword  = "password"
)

var (
	ErrInvalidExpiresInFormat = errors.New("invalid format for expires_in")
	ErrInvalidExpiresAtFormat = errors.New("invalid format for expires_at")
	ErrInvalidValidFromFormat = errors.New("invalid valid_from format")
	ErrInvalidMaxHitsFormat   = errors.New("invalid max_hits format")
)

// NewFromForm returns a new Link using values from an HTTP form.
func NewFromForm(form url.Values) (*Link, error) {
	now := time.Now()

	lnk := &Link{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
		HitCount:  0,
	}

	if value := strings.TrimSpace(form.Get(fieldExpiresIn)); value != "" {
		days, err := strconv.Atoi(value)
		if err != nil {
			return nil, ErrInvalidExpiresInFormat
		}
		lnk.ExpiresAt = now.Add(time.Hour * 24 * time.Duration(days))
	} else if value := strings.TrimSpace(form.Get(fieldExpiresAt)); value != "" {
		expiresAt, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, ErrInvalidExpiresAtFormat
		}
		lnk.ExpiresAt = expiresAt
	} else {
		return nil, ErrMissingExpiration
	}

	if value := strings.TrimSpace(form.Get(fieldValidFrom)); value != "" {
		validFrom, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, ErrInvalidValidFromFormat
		}
		lnk.ValidFrom = &validFrom
	}

	if value := strings.TrimSpace(form.Get(fieldMaxHits)); value != "" {
		maxHits, err := strconv.Atoi(value)
		if err != nil {
			return nil, ErrInvalidMaxHitsFormat
		}
		lnk.MaxHits = &maxHits
	}

	lnk.Target = strings.TrimSpace(form.Get(fieldTarget))

	if err := ValidateNewLink(lnk); err != nil {
		return nil, err
	}

	// TODO
	// lnk.Slug = generateSlug()
	// lnk.AdminToken = generateAdminToken()

	return lnk, nil
}
