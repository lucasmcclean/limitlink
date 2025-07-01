package link

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidURLFormat        = errors.New("invalid URL format")
	ErrURLSchemeNotHTTPorHTTPS = errors.New("URL must start with http or https")
	ErrURLMissingHost          = errors.New("URL must include a valid host")

	ErrExpiresAtTooSoon = errors.New("expiration time must be at least 1 minute from now")
	ErrExpiresAtTooFar  = errors.New("expiration time must be within the next 30 days")

	ErrMaxHitsNegative = errors.New("max number of hits must be 0 or greater")

	ErrValidFromTooSoon = errors.New("start time must be at least 1 minute from now")
	ErrValidFromTooFar  = errors.New("start time must be within the next 30 days")

	ErrValidFromAfterExpiresAt   = errors.New("start time must be before expiration time")
	ErrAdminExpiresBeforeExpires = errors.New("admin expiration must be after normal expiration")
)

// validateTarget ensures the target string is a valid HTTP/HTTPS URL with a host.
func validateTarget(target string) error {
	parsed, err := url.ParseRequestURI(target)
	if err != nil {
		return ErrInvalidURLFormat
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "https" && scheme != "http" {
		return ErrURLSchemeNotHTTPorHTTPS
	}
	if parsed.Host == "" {
		return ErrURLMissingHost
	}
	return nil
}

// validateExpiresAt checks that the expiration time is at least 1 minute in the future,
// but no more than 30 days ahead from the reference time 'now'.
func validateExpiresAt(expiresAt time.Time, now time.Time) error {
	if expiresAt.Before(now.Add(time.Minute)) {
		return ErrExpiresAtTooSoon
	}
	if expiresAt.After(now.Add(24 * time.Hour * 30)) {
		return ErrExpiresAtTooFar
	}
	return nil
}

// validateMaxHits verifies that maxHits is non-negative if specified (nil means no limit).
func validateMaxHits(maxHits *int) error {
	if maxHits == nil {
		return nil
	}
	if *maxHits < 0 {
		return ErrMaxHitsNegative
	}
	return nil
}

// validateValidFrom ensures the start time is at least 1 minute in the future and
// no more than 30 days ahead from 'now'.
func validateValidFrom(validFrom *time.Time, now time.Time) error {
	if validFrom == nil {
		return nil
	}
	if validFrom.Before(now.Add(time.Minute)) {
		return ErrValidFromTooSoon
	}
	if validFrom.After(now.Add(24 * time.Hour * 30)) {
		return ErrValidFromTooFar
	}
	return nil
}

// validateTimes checks cross-field logic between ValidFrom, ExpiresAt, and AdminExpiresAt.
func validateTimes(validFrom *time.Time, expiresAt, adminExpiresAt time.Time) error {
	if validFrom != nil && validFrom.After(expiresAt) {
		return ErrValidFromAfterExpiresAt
	}
	if adminExpiresAt.Before(expiresAt) {
		return ErrAdminExpiresBeforeExpires
	}
	return nil
}
