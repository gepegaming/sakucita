package service

type MidtransNotification struct {
	Currency          string `json:"currency"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	TransactionTime   string `json:"transaction_time"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	ExpiryTime        string `json:"expiry_time,omitempty"`

	CustomerDetails struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	} `json:"customer_details"`
}
