package domain

import "time"

type PaymentChannel struct {
	ID                       int32
	Code                     string
	Name                     string
	GatewayFeeFixed          int64
	GatewayFeePercentageBps  int32
	PlatformFeeFixed         int64
	PlatformFeePercentageBps int32
	IsActive                 bool
	CreatedAt                time.Time
}
