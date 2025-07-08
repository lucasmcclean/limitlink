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

	ErrAdminExpiresRequiresExpires = errors.New("admin_expires_at can only be updated if expires_at is also updated")
	ErrUpdatedAtNotSet             = errors.New("updated_at timestamp must be set")
)

// Validated represents a Link that has been successfully validated.
//
// It can only be created through the Validate function to ensure all validation
// rules have been passed. Use the Link method to access the underlying Link.
type Validated struct {
	link *Link
}

// Link returns the underlying validated Link instance.
func (v *Validated) Link() *Link {
	return v.link
}

// SetSlug applies the provided slug to the underlying Link instance's Slug.
func (v *Validated) SetSlug(slug string) {
	v.link.Slug = slug
}

// SetAdminToken applies the provided token to the underlying Link instance's
// AdminToken.
func (v *Validated) SetAdminToken(token string) {
	v.link.AdminToken = token
}

// Validate checks the provided Link's fields for correctness and consistency.
//
// It runs individual validations on fields such as Target, ExpiresAt, MaxHits,
// ValidFrom, and performs cross-field validation on ValidFrom, ExpiresAt, and
// AdminExpiresAt. If all validations pass, it returns a Validated wrapper containing
// the original Link. Otherwise, it returns an error indicating the validation failure.
func Validate(link *Link, now time.Time) (*Validated, error) {
	if err := validateTarget(link.Target); err != nil {
		return nil, err
	}
	if err := validateExpiresAt(link.ExpiresAt, now); err != nil {
		return nil, err
	}
	if err := validateMaxHits(link.MaxHits); err != nil {
		return nil, err
	}
	if err := validateValidFrom(link.ValidFrom, now); err != nil {
		return nil, err
	}
	if err := validateTimes(link.ValidFrom, link.ExpiresAt, link.AdminExpiresAt); err != nil {
		return nil, err
	}

	return &Validated{link: link}, nil
}

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

// ValidatedPatch represents a PatchLink that has been successfully validated.
//
// It can only be created through the ValidatePatch function to ensure all
// validation rules have been passed. Use the Patch method to access the
// underlying PatchLink.
type ValidatedPatch struct {
	patch *PatchLink
}

// Patch returns the underlying validated PatchLink.
func (v *ValidatedPatch) Patch() *PatchLink {
	return v.patch
}

// ValidatePatch validates updates from patch against the original Link state,
// ensuring all fields follow the rules and cross-field dependencies.
// Returns a ValidatedPatch on success, or an error on validation failure.
func ValidatePatch(original *Link, patch *PatchLink, now time.Time) (*ValidatedPatch, error) {
	if !patch.MaxHits.Remove && patch.MaxHits.Value != nil {
		if err := validateMaxHits(patch.MaxHits.Value); err != nil {
			return nil, err
		}
	}

	var validFrom *time.Time
	if patch.ValidFrom.Remove {
		validFrom = nil
	} else if patch.ValidFrom.Value != nil {
		validFrom = patch.ValidFrom.Value
	} else {
		validFrom = original.ValidFrom
	}

	if err := validateValidFrom(validFrom, now); err != nil {
		return nil, err
	}

	var expiresAt time.Time
	if patch.ExpiresAt != nil {
		expiresAt = *patch.ExpiresAt
	} else {
		expiresAt = original.ExpiresAt
	}

	if err := validateExpiresAt(expiresAt, now); err != nil {
		return nil, err
	}

	var adminExpiresAt time.Time
	if patch.AdminExpiresAt != nil {
		adminExpiresAt = *patch.AdminExpiresAt
	} else {
		adminExpiresAt = original.AdminExpiresAt
	}

	if patch.AdminExpiresAt != nil && patch.ExpiresAt == nil {
		return nil, ErrAdminExpiresRequiresExpires
	}

	if err := validateTimes(validFrom, expiresAt, adminExpiresAt); err != nil {
		return nil, err
	}

	if patch.UpdatedAt.IsZero() {
		return nil, ErrUpdatedAtNotSet
	}

	return &ValidatedPatch{patch: patch}, nil
}
