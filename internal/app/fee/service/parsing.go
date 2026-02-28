package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
)

func userFeeRowRepoToDomain(r repository.GetUserFeeRow) domain.UserFee {
	return domain.UserFee{
		PaymentChannelID:         r.PaymentChannelID,
		PlatformFeeFixed:         r.PlatformFeeFixed,
		PlatformFeePercentageBps: r.PlatformFeePercentageBps,
		GatewayFeeFixed:          r.GatewayFeeFixed,
		GatewayFeePercentageBps:  r.GatewayFeePercentageBps,
	}
}
