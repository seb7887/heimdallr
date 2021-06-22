package health

import (
	"context"

	"github.com/seb7887/heimdallr/storage"
)

type Service interface {
	GetHealth(ctx context.Context) error
}

type service struct {
	repository storage.Repository
}

func NewService(repo storage.Repository) Service {
	return &service{repository: repo}
}

func (s *service) GetHealth(ctx context.Context) error {
	return s.repository.Health(ctx)
}
