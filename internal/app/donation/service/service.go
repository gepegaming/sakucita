package service

import (
	"sakucita/internal/database/repository"
	"sakucita/internal/domain"

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
