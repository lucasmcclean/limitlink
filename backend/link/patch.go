package link

import (
	"errors"
	"time"
)

// Field is a wrapper type that allows distinguishing between an update, a no-op, and a deletion.
type Field[T any] struct {
	Value  *T
	Remove bool
}

// PatchLink represents a partial update to an existing Link.
// Use the Field type to signal if a field should be updated or explicitly removed.
type PatchLink struct {
	MaxHits        Field[int]       `bson:"-"`                          // Optional: set or remove max hit count
	PasswordHash   Field[string]    `bson:"-"`                          // Optional: set or remove password hash
	ValidFrom      Field[time.Time] `bson:"-"`                          // Optional: set or remove start time
	ExpiresAt      *time.Time       `bson:"expires_at,omitempty"`       // New expiration timestamp (or nil to skip)
	AdminExpiresAt *time.Time       `bson:"admin_expires_at,omitempty"` // New expiration timestamp (or nil to skip)
	UpdatedAt      time.Time        `bson:"updated_at"`                 // Always updated timestamp
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
		return nil, errors.New("admin_expires_at can only be updated if expires_at is also updated")
	}

	if err := validateTimes(validFrom, expiresAt, adminExpiresAt); err != nil {
		return nil, err
	}

	if patch.UpdatedAt.IsZero() {
		return nil, errors.New("updated_at timestamp must be set")
	}

	return &ValidatedPatch{patch: patch}, nil
}
