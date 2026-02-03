package midtrans

import (
	"time"

	"sakucita/pkg/config"

	"resty.dev/v3"
)

type MidtransClient struct {
	http *resty.Client
}

func NewMidtransClient(config config.App) *MidtransClient {
	c := resty.New().
		SetBaseURL(config.Midtrans.BaseURL).
		SetBasicAuth(config.Midtrans.ServerKey, "").
		SetTimeout(15*time.Second).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	return &MidtransClient{http: c}
}
