package link

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Link represents a shortened URL with optional access controls and usage limits.
//
// It maps a short slug to a target URL, supporting features such as expiration,
// max hit count, valid time windows, password protection, and admin management.
type Link struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`           // MongoDB ID
	Slug          string             `bson:"slug"`                    // Unique identifier for the link
	Target        string             `bson:"target"`                  // Destination URL
	CreatedAt     time.Time          `bson:"created_at"`              // Creation timestamp
	ExpiresAt     time.Time          `bson:"expires_at"`              // Expiration timestamp
	UpdatedAt     time.Time          `bson:"updated_at"`              // Last updated timestamp
	MaxHits       *int               `bson:"max_hits,omitempty"`      // Optional max allowed hits
	HitCount      int                `bson:"hit_count"`               // Number of hits so far
	ValidFrom     *time.Time         `bson:"valid_from,omitempty"`    // Optional start validity timestamp
	AdminToken    string             `bson:"admin_token"`             // Owner’s admin token
	PasswordHash  *string            `bson:"password_hash,omitempty"` // Optional password hash for access
	SchemaVersion int                `bson:"schema_version"`          // Schema version for migration
}
