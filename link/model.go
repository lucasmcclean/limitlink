package link

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Link represents a shortened URL with optional access controls, expiration,
// and usage limits.
//
// Each Link maps a short slug to a target URL. It can enforce access rules such
// as expiration time, maximum hit count, time windows, and optional password
// protection. Links can also be managed via an associated admin token.
type Link struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Slug          string             `bson:"slug"`
	Target        string             `bson:"target"`
	CreatedAt     time.Time          `bson:"created_at"`
	ExpiresAt     time.Time          `bson:"expires_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	MaxHits       *int               `bson:"max_hits,omitempty"`
	HitCount      int                `bson:"hit_count"`
	ValidFrom     *time.Time         `bson:"valid_from,omitempty"`
	AdminToken    string             `bson:"admin_token"`
	PasswordHash  *string            `bson:"password_hash,omitempty"`
	SchemaVersion int                `bson:"schema_version"`
}
