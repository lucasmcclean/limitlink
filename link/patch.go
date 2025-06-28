package link

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// Extended set of field names for removal form input keys.
const (
	fieldRemoveMaxHits   = "remove-max-hits"
	fieldRemoveValidFrom = "remove-valid-from"
	fieldRemovePassword  = "remove-password"
)

// Field is a wrapper type that allows distinguishing between an update, a no-op, and a deletion.
type Field[T any] struct {
	Value  *T
	Remove bool
}

// PatchLink represents a partial update to an existing Link.
// Use the Field type to signal if a field should be updated or explicitly removed.
type PatchLink struct {
	ExpiresAt    *time.Time       `bson:"expires_at,omitempty"` // New expiration timestamp (or nil to skip)
	UpdatedAt    time.Time        `bson:"updated_at"`           // Always updated timestamp
	MaxHits      Field[int]       `bson:"-"`                    // Optional: set or remove max hit count
	ValidFrom    Field[time.Time] `bson:"-"`                    // Optional: set or remove start time
	PasswordHash Field[string]    `bson:"-"`                    // Optional: set or remove password hash
}

// NewPatchLink initializes a PatchLink with the current time as UpdatedAt.
func NewPatchLink() *PatchLink {
	return &PatchLink{
		UpdatedAt: time.Now(),
	}
}

// PatchFromForm builds a PatchLink from form values for partial updates.
// It respects "remove-*" checkboxes by setting Field.Remove = true.
func PatchFromForm(form url.Values) (*PatchLink, error) {
	now := time.Now()
	patch := NewPatchLink()

	expiresAtStr := getValue(fieldExpiresAt, form)
	if expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return nil, errors.New("Expiration date must be in a valid format (e.g., 2025-12-31T23:59:59Z).")
		}
		if err := validateExpiresAt(expiresAt, now); err != nil {
			return nil, err
		}
		patch.ExpiresAt = &expiresAt
	}

	if form.Has(fieldRemoveMaxHits) {
		patch.MaxHits.Remove = true
	} else {
		maxHitsStr := getValue(fieldMaxHits, form)
		if maxHitsStr != "" {
			maxHits, err := strconv.Atoi(maxHitsStr)
			if err != nil || maxHits < 1 {
				return nil, errors.New("Max hits must be a whole number greater than 0.")
			}
			if err := validateMaxHits(&maxHits); err != nil {
				return nil, err
			}
			patch.MaxHits.Value = &maxHits
		}
	}

	if form.Has(fieldRemoveValidFrom) {
		patch.ValidFrom.Remove = true
	} else {
		validFromStr := getValue(fieldValidFrom, form)
		if validFromStr != "" {
			validFrom, err := time.Parse(time.RFC3339, validFromStr)
			if err != nil {
				return nil, errors.New("Start date must be in a valid format (e.g., 2025-06-01T00:00:00Z).")
			}
			if err := validateValidFrom(&validFrom, now); err != nil {
				return nil, err
			}
			if patch.ExpiresAt != nil {
				if err := validateValidFromWithExpiresAt(&validFrom, *patch.ExpiresAt); err != nil {
					return nil, err
				}
			}
			patch.ValidFrom.Value = &validFrom
		}
	}

	if form.Has(fieldRemovePassword) {
		patch.PasswordHash.Remove = true
	} else {
		password := getValue(fieldPassword, form)
		if password != "" {
			hash, err := generateHash(password)
			if err != nil {
				return nil, errors.New("Error generating the password hash, please try again.")
			}
			patch.PasswordHash.Value = &hash
		}
	}

	return patch, nil
}
