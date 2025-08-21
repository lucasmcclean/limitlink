package link

import "time"

// PublicLink is a safe-to-share representation of a Link.
type PublicLink struct {
	Slug           string     `bson:"slug" json:"slug"`                                // Unique identifier for the link
	AdminToken     string     `bson:"admin_token" json:"adminToken"`                   // Owner’s admin token
	Target         string     `bson:"target" json:"target"`                            // Destination URL
	HitCount       int        `bson:"hit_count" json:"hitCount"`                       // Number of hits so far
	MaxHits        *int       `bson:"max_hits,omitempty" json:"maxHits,omitempty"`     // Optional max allowed hits
	ValidFrom      *time.Time `bson:"valid_from,omitempty" json:"validFrom,omitempty"` // Optional start validity timestamp
	CreatedAt      time.Time  `bson:"created_at" json:"createdAt"`                     // Creation timestamp
	ExpiresAt      time.Time  `bson:"expires_at" json:"expiresAt"`                     // Expiration timestamp
	AdminExpiresAt time.Time  `bson:"admin_expires_at" json:"adminExpiresAt"`          // Expiration timestamp for admin access
	UpdatedAt      time.Time  `bson:"updated_at" json:"updatedAt"`                     // Last updated timestamp
}

func (lnk *Link) ToPublic() *PublicLink {
	return &PublicLink{
		Slug:           lnk.Slug,
		AdminToken:     lnk.AdminToken,
		Target:         lnk.Target,
		HitCount:       lnk.HitCount,
		MaxHits:        lnk.MaxHits,
		ValidFrom:      lnk.ValidFrom,
		CreatedAt:      lnk.CreatedAt,
		ExpiresAt:      lnk.ExpiresAt,
		AdminExpiresAt: lnk.AdminExpiresAt,
		UpdatedAt:      lnk.UpdatedAt,
	}
}
