package service

import "sakucita/internal/infra/postgres/repository"

func mapMidtransStatus(status string, fraud string) repository.TransactionStatus {
	switch status {
	case "settlement", "capture":
		if fraud == "" || fraud == "accept" {
			return repository.TransactionStatusPAID
		}
		return repository.TransactionStatusFAILED

	case "pending":
		return repository.TransactionStatusPENDING

	case "expire":
		return repository.TransactionStatusEXPIRED

	case "cancel", "deny":
		return repository.TransactionStatusFAILED

	default:
		return repository.TransactionStatusFAILED
	}
}
