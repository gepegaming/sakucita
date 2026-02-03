package midtrans

import (
	"context"
	"fmt"
)

func (c *MidtransClient) CreateQRIS(
	ctx context.Context,
	req MidtransQRISRequest,
) (*MidtransQRISResponse, error) {
	var resp MidtransQRISResponse

	r, err := c.http.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&resp).
		Post("/v2/charge")
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("midtrans error: %s", r.String())
	}

	return &resp, nil
}
