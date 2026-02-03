package midtrans

type QRISAcquirer string

const (
	QRISGopay     QRISAcquirer = "gopay"
	QRISShopeePay QRISAcquirer = "airpay shopee"
)

type MidtransQRISRequest struct {
	PaymentType string `json:"payment_type"`

	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	ItemDetails        []MidtransItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    *MidtransCustomerDetails   `json:"customer_details,omitempty"`

	QRIS MidtransQRISDetail `json:"qris"`
}
type MidtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type MidtransItemDetail struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int32  `json:"quantity"`
	Name     string `json:"name"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type MidtransQRISDetail struct {
	Acquirer string `json:"acquirer"`
}

type MidtransQRISResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	Acquirer          string `json:"acquirer"`
	QRString          string `json:"qr_string"`
	ExpiryTime        string `json:"expiry_time"`
}
