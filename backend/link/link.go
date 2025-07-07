package link

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// adminTokenLen is the number of characters in a generated admin token.
const adminTokenLen = 22

// schemaVersion defines the current version of the link schema.
const schemaVersion = 1

// Link represents a shortened URL with optional access controls and usage
// limits.
//
// It maps a short slug to a target URL, supporting features such as expiration,
// max hit count, valid time windows, password protection, and admin management.
type Link struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`           // MongoDB ID
	Slug           string             `bson:"slug"`                    // Unique identifier for the link
	AdminToken     string             `bson:"admin_token"`             // Owner’s admin token
	Target         string             `bson:"target"`                  // Destination URL
	HitCount       int                `bson:"hit_count"`               // Number of hits so far
	MaxHits        *int               `bson:"max_hits,omitempty"`      // Optional max allowed hits
	PasswordHash   *string            `bson:"password_hash,omitempty"` // Optional password hash for access
	ValidFrom      *time.Time         `bson:"valid_from,omitempty"`    // Optional start validity timestamp
	CreatedAt      time.Time          `bson:"created_at"`              // Creation timestamp
	ExpiresAt      time.Time          `bson:"expires_at"`              // Expiration timestamp
	AdminExpiresAt time.Time          `bson:"admin_expires_at"`        // Expiration timestamp for admin access
	UpdatedAt      time.Time          `bson:"updated_at"`              // Last updated timestamp
	SchemaVersion  int                `bson:"schema_version"`          // Schema version for migration
}

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
