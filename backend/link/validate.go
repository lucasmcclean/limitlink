package link

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidURLFormat        = errors.New("invalid URL format")
	ErrURLSchemeNotHTTPorHTTPS = errors.New("URL must start with http or https")
	ErrURLMissingHost          = errors.New("URL must include a valid host")

	ErrExpiresAtTooSoon = fmt.Errorf("expiration time must be at least %d minute from now", minTime)
	ErrExpiresAtTooFar  = fmt.Errorf("expiration time must be within the next %d days", maxTime)

	ErrMaxHitsNegative = errors.New("max number of hits must be 0 or greater")
	ErrMaxHitsTooLarge = fmt.Errorf("max number of hits must be less than %d", maxMaxHits)

	ErrPasswordTooLong = fmt.Errorf("password is too long (max %d characters)", maxPasswordLen)

	ErrValidFromTooSoon = fmt.Errorf("start time must be at least %d minute from now", minTime)
	ErrValidFromTooFar  = fmt.Errorf("start time must be within the next %d days", maxTime)

	ErrValidFromAfterExpiresAt   = errors.New("start time must be before expiration time")
	ErrAdminExpiresBeforeExpires = errors.New("admin expiration must be after normal expiration")

	ErrAdminExpiresRequiresExpires = errors.New("admin_expires_at can only be updated if expires_at is also updated")
	ErrUpdatedAtNotSet             = errors.New("updated_at timestamp must be set")

	ErrUnrecognizedCharset = errors.New("unrecognized character set")
	ErrInvalidSlugLen      = errors.New("slug length must be between 6 and 12 inclusive")
	ErrGeneratingSlug      = errors.New("error generating the slug")

	ErrHashingPassword = errors.New("error hashing password")
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

// SetSlug generates and applies a validated slug to the underlying link.
func (v *Validated) SetSlug(length int, charset string) error {
	if length < MinSlugLen || length > MaxSlugLen {
		return ErrInvalidSlugLen
	}

	slug, err := generateSlug(length, charset)
	if err != nil {
		return err
	}

	v.link.Slug = slug
	return nil
}

// SetAdminToken generates and applies a validated token to the underlying link.
func (v *Validated) SetAdminToken() error {
	token, err := generateAdminToken(adminTokenLen)
	if err != nil {
		return err
	}
	v.link.AdminToken = token
	return nil
}

// SetPasswordHash generates and applies a password hash to the underlying link
// once the password has been validated.
func (v *Validated) SetPasswordHash(password *string) error {
	if password == nil || *password == "" {
		return nil
	}

	if len(*password) > maxPasswordLen {
		return ErrPasswordTooLong
	}

	passwordHash, err := hashPassword(*password)
	if err != nil {
		return ErrHashingPassword
	}

	v.link.PasswordHash = &passwordHash
	return nil
}

// Validate checks the provided Link's fields for correctness and consistency.
//
// Do not assign a Slug, AdminToken, or Password before validating.
// Use the provided SetSlug, SetAdminToken, and SetPasswordHash functions.
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

// SetPasswordHash generates and applies a password hash to the underlying link
// once the password has been validated.
func (v *ValidatedPatch) SetPasswordHash(password *string) error {
	if password == nil || *password == "" {
		return nil
	}

	if len(*password) > maxPasswordLen {
		return ErrPasswordTooLong
	}

	passwordHash, err := hashPassword(*password)
	if err != nil {
		return ErrHashingPassword
	}

	v.patch.PasswordHash = Field[string]{
		Value:  &passwordHash,
		Remove: false,
	}
	return nil
}

// ValidatePatch validates updates from patch against the original Link state,
// ensuring all fields follow the rules and cross-field dependencies.
//
// Do not assign a Password before validating; use the provided SetPasswordHash
// instead.
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

	if err := validateValidFrom(validFrom, now); err != nil {
		return nil, err
	}

	return &ValidatedPatch{patch: patch}, nil
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
	if expiresAt.Before(now.Add(minTime)) {
		return ErrExpiresAtTooSoon
	}
	if expiresAt.After(now.Add(maxTime)) {
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
	if validFrom.Before(now.Add(minTime)) {
		return ErrValidFromTooSoon
	}
	if validFrom.After(now.Add(maxTime)) {
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
	} else if *maxHits > maxMaxHits {
		return ErrMaxHitsTooLarge
	}
	return nil
}
