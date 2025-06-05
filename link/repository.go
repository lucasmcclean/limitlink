package link

import "context"

type Repository interface {
	Create(ctx context.Context, link *Link) error
	GetBySlug(ctx context.Context, slug string) (*Link, error)
	IncBySlug(ctx context.Context, slug string) error
	GetByToken(ctx context.Context, adminToken string) (*Link, error)
	DeleteByToken(ctx context.Context, adminToken string) error
	UpdateByToken(ctx context.Context, adminToken string, updated *Link) error
}
