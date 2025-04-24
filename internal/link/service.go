package link

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Link represents a shortened URL and its constraints
type Link struct {
	ID         uuid.UUID  `db:"id"`
	Original   string     `db:"original"`
	Short      string     `db:"short"`
	AdminToken uuid.UUID  `db:"admin_token"`
	MaxUses    *int       `db:"max_uses"`
	ClickCount int        `db:"click_count"`
	ExpiresAt  *time.Time `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

// CreateLinkRequest represents the incoming request for creating a link
type CreateLinkRequest struct {
	Original  string     `json:"original"`
	Short     string     `json:"short"`
	MaxUses   *int       `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Service encapsulates business rules for links
type Service struct {
	Repo LinkRepository
}

// NewService creates a new link service with the given repository
func NewService(repo LinkRepository) *Service {
	return &Service{Repo: repo}
}

// CreateLink handles link creation logic
func (s *Service) CreateLink(ctx context.Context, req CreateLinkRequest) (Link, error) {
	adminToken, err := s.Repo.CreateLink(req.Original, req.Short, req.MaxUses, req.ExpiresAt)
	if err != nil {
		return Link{}, err
	}

	link, err := s.Repo.GetByAdminToken(ctx, adminToken)
	if err != nil {
		return Link{}, err
	}

	return link, nil
}

// VisitLink resolves a shortened link to the original URL
func (s *Service) VisitLink(ctx context.Context, short string) (string, error) {
	link, err := s.Repo.GetByShort(ctx, short)
	if err != nil {
		return "", err
	}

	return link.Original, nil
}

// GetByAdminToken fetches a link by its admin token
func (s *Service) GetByAdminToken(ctx context.Context, token uuid.UUID) (Link, error) {
	return s.Repo.GetByAdminToken(ctx, token)
}

// DeleteByAdminToken deletes a link by its admin token
func (s *Service) DeleteByAdminToken(ctx context.Context, token uuid.UUID) error {
	return s.Repo.DeleteByAdminToken(ctx, token)
}

// DeleteExpired removes expired links
func (s *Service) DeleteExpired(ctx context.Context) error {
	return s.Repo.DeleteExpired(ctx)
}
