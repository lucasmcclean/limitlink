package link

import (
	"context"
)

// Repository defines persistence operations for Link objects.
type Repository interface {
	// Create inserts a new link into the repository.
	Create(ctx context.Context, link *Link) error

	// GetBySlug retrieves a link by its public slug.
	GetBySlug(ctx context.Context, slug string) (*Link, error)

	// IncBySlug increments the hit count for the given slug.
	IncBySlug(ctx context.Context, slug string) error

	// GetByToken retrieves a link by its admin token.
	GetByToken(ctx context.Context, token string) (*Link, error)

	// DeleteByToken removes a link by its admin token.
	DeleteByToken(ctx context.Context, token string) error

	// PatchByToken updates a link by its admin token.
	PatchByToken(ctx context.Context, token string, updated *PatchLink) error
}
