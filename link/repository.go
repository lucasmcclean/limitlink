package link

import (
	"context"
)

// Repository defines persistence operations for Link objects.
type Repository interface {
	// Close releases any held resources (e.g., DB connections).
	Close(ctx context.Context) error

	// Create inserts a new link into the repository.
	Create(ctx context.Context, link *Link) error

	// GetBySlug retrieves a link by its public slug.
	GetBySlug(ctx context.Context, slug string) (*Link, error)

	// IncBySlug atomically increments the hit count for the given slug.
	IncBySlug(ctx context.Context, slug string) error

	// GetByToken retrieves a link by its admin token (used for managing the link).
	GetByToken(ctx context.Context, adminToken string) (*Link, error)

	// DeleteByToken removes a link identified by its admin token.
	DeleteByToken(ctx context.Context, adminToken string) error

	// PatchByToken updates a link using its admin token as an identifier.
	PatchByToken(ctx context.Context, adminToken string, updated *PatchLink) error
}
