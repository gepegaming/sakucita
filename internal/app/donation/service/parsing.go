package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"
)

func donationMessageRepoToDomain(d repository.DonationMessage) domain.DonationMessage {

	mediaStartSecond, _ := utils.PgTypeInt4ToInt32(d.MediaStartSeconds)
	maxPlaySecond, _ := utils.PgTypeInt4ToInt32(d.MaxPlaySeconds)
	pricePerSecond, _ := utils.PgTypeInt8ToInt64Ptr(d.PricePerSecond)
	meta, _ := utils.JSONBytesToMap(d.Meta)

	return domain.DonationMessage{
		ID:                d.ID,
		PayeeUserID:       d.PayeeUserID,
		PayerUserID:       utils.PgTypeUUIDToUUIDPtr(d.PayerUserID),
		PayerName:         d.PayerName,
		Message:           d.Message,
		Email:             d.Email,
		MediaType:         string(d.MediaType),
		MediaUrl:          utils.PgTypeTextToStringPtr(d.MediaUrl),
		MediaStartSeconds: mediaStartSecond,
		MaxPlaySeconds:    maxPlaySecond,
		PricePerSecond:    pricePerSecond,
		Amount:            d.Amount,
		Currency:          d.Currency,
		Status:            string(d.Status),
		Meta:              meta,
		PlayedAt:          d.PlayedAt.Time,
		CreatedAt:         d.CreatedAt.Time,
	}
}
