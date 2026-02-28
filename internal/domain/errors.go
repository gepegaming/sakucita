package domain

import "errors"

type AppError struct {
	Code    int    // HTTP Status Code (401, 404, dll)
	Message string // Pesan untuk user
	Err     error  // Error asli (optional, untuk logging internal)
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotfound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrTooManyRequests     = errors.New("too many requests")

	// general
	ErrMsgInvalidRequest      = "invalid request"
	ErrMsgInternalServerError = "oops, something went wrong, please try again later"

	// user domain
	ErrMsgEmailAlreadyExists    = "email already exists"
	ErrMsgPhoneAlreadyExists    = "phone already exists"
	ErrMsgNicknameAlreadyExists = "nickname already exists"
	ErrMsgUserNotFound          = "user not found"

	// payment channel
	ErrMsgPaymentChannelNotFound = "payment channel not found"

	// auth domain
	ErrMsgInvalidCredentials = "invalid credentials"
	ErrMsgSessionNotFound    = "session not found"
	ErrMsgDeviceIdMissmatch  = "device id missmatch"
	ErrMsgUnauthorized       = "unauthorized"

	// transaction domain
	ErrMsgTransactionNotFound = "transaction not found"

	// donation domain
	ErrMsgDonationMessageNotFound = "donation message not found"
	ErrMsgDonationNotFound        = "donation not found"
	ErrFailedToCreateDonation     = "failed to create donation"
)
