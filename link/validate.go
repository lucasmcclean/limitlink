package link

import (
	"errors"
	"time"
)

var (
	ErrLinkExpired         = errors.New("link expired")
	ErrLinkNotYetValid     = errors.New("link not yet valid")
	ErrLinkExceededMaxHits = errors.New("link exceeded max hits")
)

func validate(link *Link) error {
	now := time.Now()

	if now.After(link.ExpiresAt) {
		return ErrLinkExpired
	}
	if link.ValidFrom != nil && now.Before(*link.ValidFrom) {
		return ErrLinkNotYetValid
	}
	if link.MaxHits != nil && link.HitCount >= *link.MaxHits {
		return ErrLinkExceededMaxHits
	}
	return nil
}
