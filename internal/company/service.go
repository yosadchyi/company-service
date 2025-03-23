package company

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the expected storage operations for Company entities.
type Repository interface {
	Create(ctx context.Context, c *Company) error
	Get(ctx context.Context, id uuid.UUID) (*Company, error)
	Update(ctx context.Context, c *Company) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, c *Company) error {
	return s.repo.Create(ctx, c)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Company, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, patch *PartialCompany) (*Company, error) {
	// Load existing company
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply patch fields only if set
	if patch.Name != nil {
		current.Name = *patch.Name
	}
	if patch.Description != nil {
		current.Description = patch.Description
	}
	if patch.AmountOfEmployees != nil {
		current.AmountOfEmployees = *patch.AmountOfEmployees
	}
	if patch.Registered != nil {
		current.Registered = *patch.Registered
	}
	if patch.Type != nil {
		current.Type = *patch.Type
	}

	// Persist updated company
	if err := s.repo.Update(ctx, current); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
