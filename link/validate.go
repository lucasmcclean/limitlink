package link

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

const (
	MinLifetime   = 5 * time.Minute // Minimum time a link must be valid
	MinValidDelay = 1 * time.Minute // Minimum delay before ValidFrom
)

var (
	ErrLinkExpired         = errors.New("link expired")
	ErrLinkNotYetValid     = errors.New("link not yet valid")
	ErrLinkExceededMaxHits = errors.New("link exceeded max hits")

	ErrInvalidExpiration = errors.New(fmt.Sprintf(
		"expiration must be at least %d minutes from now",
		MinLifetime))
	ErrInvalidValidFrom = errors.New(fmt.Sprintf(
		"valid_from must be before expires_at, not in the past, and be at least %d minutes from now",
		MinValidDelay))
	ErrInvalidMaxHits    = errors.New("max_hits must be greater than 0")
	ErrMissingTarget     = errors.New("target is required")
	ErrInvalidTarget     = errors.New("target must be a valid URL with scheme and host")
	ErrMissingExpiration = errors.New("expiration is required")
)

// Validate is used to ensure that a link is still valid before serving.
func Validate(l *Link) error {
	now := time.Now()

	if now.After(l.ExpiresAt) {
		return ErrLinkExpired
	}
	if l.ValidFrom != nil && now.Before(*l.ValidFrom) {
		return ErrLinkNotYetValid
	}
	if l.MaxHits != nil && l.HitCount >= *l.MaxHits {
		return ErrLinkExceededMaxHits
	}
	return nil
}

// ValidateNewLink validates a newly created link before insertion.
func ValidateNewLink(l *Link) error {
	if err := validateTarget(l.Target); err != nil {
		return err
	}
	if err := validateExpiration(l.ExpiresAt); err != nil {
		return err
	}
	if err := validateValidFrom(l.ValidFrom, l.ExpiresAt); err != nil {
		return err
	}
	if err := validateMaxHits(l.MaxHits); err != nil {
		return err
	}
	return nil
}

func validateExpiration(expiresAt time.Time) error {
	if expiresAt.IsZero() {
		return ErrMissingExpiration
	}
	if expiresAt.Before(time.Now().Add(MinLifetime)) {
		return ErrInvalidExpiration
	}
	return nil
}

func validateValidFrom(validFrom *time.Time, expiresAt time.Time) error {
	if validFrom == nil {
		return nil
	}
	now := time.Now()
	if validFrom.After(expiresAt) {
		return ErrInvalidValidFrom
	}
	if validFrom.Before(now.Add(MinValidDelay)) {
		return ErrInvalidValidFrom
	}
	return nil
}

func validateMaxHits(maxHits *int) error {
	if maxHits == nil {
		return nil
	}
	if *maxHits <= 0 {
		return ErrInvalidMaxHits
	}
	return nil
}

func validateTarget(target string) error {
	if target == "" {
		return ErrMissingTarget
	}
	parsed, err := url.ParseRequestURI(target)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrInvalidTarget
	}
	return nil
}
