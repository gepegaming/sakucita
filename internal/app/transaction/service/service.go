package service

import (
	"context"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/pkg/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type TransactionService interface {
	GetByExternalReference(ctx context.Context, ref string) (transaction, error)
	MarkAsSuccess(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	MarkAsFailed(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error
}

type transaction struct {
	q      *repository.Queries
	config config.App
	log    zerolog.Logger
}

func (t *transaction) GetByExternalReference(ctx context.Context, ref string) (transaction, error) {
	panic("unimplemented")
}

func (t *transaction) MarkAsFailed(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	panic("unimplemented")
}

func (t *transaction) MarkAsSuccess(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	panic("unimplemented")
}

func (t *transaction) UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status string) error {
	panic("unimplemented")
}

func NewService(q *repository.Queries, config config.App, log zerolog.Logger) TransactionService {
	return &transaction{q, config, log}
}
