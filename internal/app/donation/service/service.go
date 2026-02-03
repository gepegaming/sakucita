package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"

	"github.com/rs/zerolog"
)

type service struct {
	q   *repository.Queries
	log zerolog.Logger
}

func NewService(
	q *repository.Queries,
	log zerolog.Logger,
) domain.DonationService {
	return &service{q, log}
}

func (s *service) CreateDonation() {
	panic("unimplemented")
}
