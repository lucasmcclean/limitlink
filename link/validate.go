package link

import (
	"errors"
	"net/url"
	"time"
)

// validateTarget ensures the target string is a valid HTTP/HTTPS URL with a host.
func validateTarget(target string) error {
	parsed, err := url.ParseRequestURI(target)
	if err != nil {
		return errors.New("Please enter a valid URL.")
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return errors.New("URL must start with http or https.")
	}
	if parsed.Host == "" {
		return errors.New("URL must include a valid host.")
	}
	return nil
}

// validateExpiresAt checks that the expiration time is at least 1 minute in the future,
// but no more than 30 days ahead from the reference time 'now'.
func validateExpiresAt(expiresAt time.Time, now time.Time) error {
	if expiresAt.Before(now.Add(time.Minute)) {
		return errors.New("Expiration time must be at least 1 minute from now.")
	}
	if expiresAt.After(now.Add(time.Hour * 24 * 30)) {
		return errors.New("Expiration time must be within the next 30 days.")
	}
	return nil
}

// validateMaxHits verifies that maxHits is non-negative if specified (nil means no limit).
func validateMaxHits(maxHits *int) error {
	if maxHits == nil {
		return nil
	}
	if *maxHits < 0 {
		return errors.New("Max number of hits must be 0 or greater.")
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
		return errors.New("Start time must be at least 1 minute from now.")
	}
	if validFrom.After(now.Add(time.Hour * 24 * 30)) {
		return errors.New("Start time must be within the next 30 days.")
	}
	return nil
}

// validateValidFromWithExpiresAt ensures the start time precedes the expiration time.
func validateValidFromWithExpiresAt(validFrom *time.Time, expiresAt time.Time) error {
	if validFrom == nil {
		return nil
	}
	if validFrom.After(expiresAt) {
		return errors.New("Start time must be before the expiration time.")
	}
	return nil
}
