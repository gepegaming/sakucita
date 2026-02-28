package domain

type UserFee struct {
	PaymentChannelID         int32
	PlatformFeeFixed         int64
	PlatformFeePercentageBps int32
	GatewayFeeFixed          int64
	GatewayFeePercentageBps  int32
}
