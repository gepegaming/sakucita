package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
)

func paymentChannelRepoToDomain(p repository.PaymentChannel) domain.PaymentChannel {
	return domain.PaymentChannel{
		ID:                       p.ID,
		Code:                     p.Code,
		Name:                     p.Name,
		GatewayFeeFixed:          p.GatewayFeeFixed,
		GatewayFeePercentageBps:  p.GatewayFeePercentageBps,
		PlatformFeeFixed:         p.PlatformFeeFixed,
		PlatformFeePercentageBps: p.PlatformFeePercentageBps,
		IsActive:                 p.IsActive,
		CreatedAt:                p.CreatedAt.Time,
	}
}
