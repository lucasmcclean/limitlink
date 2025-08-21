package link

import (
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

// NewPatchLink initializes a PatchLink with the current time as UpdatedAt.
func NewPatchLink(now time.Time) *PatchLink {
	return &PatchLink{
		UpdatedAt: now,
	}
}
