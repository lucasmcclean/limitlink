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

const (
	fieldTarget      = "target"
	fieldExpiresIn   = "expires_in"
	fieldExpiresAt   = "expires_at"
	fieldMaxHits     = "max_hits"
	fieldValidFrom   = "valid_from"
	fieldPassword    = "password"
	fieldSlugLength  = "slug_length"
	fieldSlugCharset = "slug_charset"
)

var (
	ErrInvalidExpiresInFormat = errors.New("invalid format for expires_in")
	ErrInvalidExpiresAtFormat = errors.New("invalid format for expires_at")
	ErrInvalidValidFromFormat = errors.New("invalid valid_from format")
	ErrInvalidMaxHitsFormat   = errors.New("invalid max_hits format")
	ErrInvalidSlugCharset     = errors.New("invalid slug_charset value")
	ErrInvalidSlugLength      = errors.New(
		fmt.Sprintf("slug_length must be between %d and %d inclusive", minSlugLen, maxSlugLen),
	)
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

	slugLength := defaultSlugLen
	if val := form.Get(fieldSlugLength); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed >= 6 && parsed <= 12 {
			slugLength = parsed
		} else {
			return nil, ErrInvalidSlugLength
		}
	}

	var charset slugCharset
	switch form.Get(fieldSlugCharset) {
	case "", "alphanumeric":
		charset = alphanumeric
	case "letters":
		charset = letters
	case "numbers":
		charset = numbers
	default:
		return nil, ErrInvalidSlugCharset
	}

	if err := ValidateNewLink(lnk); err != nil {
		return nil, err
	}

	var err error
	lnk.Slug, err = generateSlug(slugLength, charset)
	if err != nil {
		return nil, err
	}

	lnk.AdminToken, err = generateAdminToken()
	if err != nil {
		return nil, err
	}

	return lnk, nil
}
