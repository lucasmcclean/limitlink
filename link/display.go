package link

import "time"

type DisplayLink struct {
	Slug       string  `bson:"slug"`
	Target     string  `bson:"target"`
	CreatedAt  string  `bson:"created_at"`
	ExpiresAt  string  `bson:"expires_at"`
	UpdatedAt  string  `bson:"updated_at"`
	MaxHits    *int    `bson:"max_hits,omitempty"`
	HitCount   int     `bson:"hit_count"`
	ValidFrom  *string `bson:"valid_from,omitempty"`
	AdminToken string  `bson:"admin_token"`
}

const isoLayout = time.RFC3339

func (lnk *Link) IntoDisplay() *DisplayLink {
	dlnk := &DisplayLink{
		Slug:       lnk.Slug,
		Target:     lnk.Target,
		CreatedAt:  lnk.CreatedAt.Format(isoLayout),
		ExpiresAt:  lnk.ExpiresAt.Format(isoLayout),
		UpdatedAt:  lnk.UpdatedAt.Format(isoLayout),
		MaxHits:    lnk.MaxHits,
		HitCount:   lnk.HitCount,
		AdminToken: lnk.AdminToken,
	}

	if lnk.ValidFrom != nil {
		formatted := lnk.ValidFrom.Format(isoLayout)
		dlnk.ValidFrom = &formatted
	}

	return dlnk
}
