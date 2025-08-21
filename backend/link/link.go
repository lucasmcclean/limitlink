package link

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// minSlugLen is the minimum length of a generated slug.
	MinSlugLen = 6

	// maxSlugLen is the maximum length of a generated slug.
	MaxSlugLen = 12

	// maxMaxHits is the maximum valid amount for max hits.
	maxMaxHits = 1_000_000

	// adminTokenLen is the number of characters in a generated admin token.
	adminTokenLen = 22

	// maxPasswordLen is the maximum allowed length of a supplied password.
	maxPasswordLen = 256

	// minTime is the minimum amount a given time must be from now to be valid.
	minTime = time.Minute

	// maxTime is the maximum amount of time in the future a provided time can be.
	maxTime = time.Hour * 24 * 30

	// schemaVersion defines the current version of the link schema.
	schemaVersion = 1
)

// Link represents a shortened URL with optional access controls and usage
// limits.
type Link struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`                          // MongoDB ID
	Slug           string             `bson:"slug" json:"slug"`                                // Unique identifier for the link
	AdminToken     string             `bson:"admin_token" json:"adminToken"`                   // Owner’s admin token
	Target         string             `bson:"target" json:"target"`                            // Destination URL
	HitCount       int                `bson:"hit_count" json:"hitCount"`                       // Number of hits so far
	MaxHits        *int               `bson:"max_hits,omitempty" json:"maxHits,omitempty"`     // Optional max allowed hits
	PasswordHash   *string            `bson:"password_hash,omitempty" json:"-"`                // Optional password hash (not exposed in JSON)
	ValidFrom      *time.Time         `bson:"valid_from,omitempty" json:"validFrom,omitempty"` // Optional start validity timestamp
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`                     // Creation timestamp
	ExpiresAt      time.Time          `bson:"expires_at" json:"expiresAt"`                     // Expiration timestamp
	AdminExpiresAt time.Time          `bson:"admin_expires_at" json:"adminExpiresAt"`          // Expiration timestamp for admin access
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`                     // Last updated timestamp
	SchemaVersion  int                `bson:"schema_version" json:"-"`                         // Schema version for migration
}

func (l *Link) IsAvailable(now time.Time) bool {
	if l.MaxHits != nil && l.HitCount >= *l.MaxHits {
		return false
	}
	if l.ValidFrom != nil && now.Before(*l.ValidFrom) {
		return false
	}
	if now.After(l.ExpiresAt) {
		return false
	}
	return true
}
