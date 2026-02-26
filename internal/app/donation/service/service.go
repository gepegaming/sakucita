package service

import (
	"context"

	"sakucita/internal/domain"
	"sakucita/internal/infra/midtrans"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type service struct {
	db             *pgxpool.Pool
	q              *repository.Queries
	log            zerolog.Logger
	midtransClient midtrans.MidtransClient
}

type DonationService interface {
	CreateDonation(ctx context.Context, req CreateDonationCommand) (*CreateDonationResult, error)
}

func NewService(
	db *pgxpool.Pool,
	q *repository.Queries,
	log zerolog.Logger,
	midtransClient midtrans.MidtransClient,
) DonationService {
	return &service{db, q, log, midtransClient}
}

func (s *service) CreateDonation(
	ctx context.Context,
	req CreateDonationCommand,
) (*CreateDonationResult, error) {

	/**
	 * 1. Setup Database Transaction
	 */
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	/**
	 * 2. Validate & Load Required Data
	 */

	creator, err := qtx.GetUserByID(ctx, req.PayeeUserID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, domain.NewAppError(
				fiber.StatusNotFound,
				domain.ErrMsgUserNotFound,
				domain.ErrNotfound,
			)
		}
		s.log.Err(err).Msg("failed to get user by id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	paymentChannel, err := qtx.GetPaymentChannelByCode(ctx, req.PaymentChannel)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, domain.NewAppError(
				fiber.StatusNotFound,
				domain.ErrMsgPaymentChannelNotFound,
				domain.ErrNotfound,
			)
		}
		s.log.Err(err).Msg("failed to get payment channel")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	fees, err := qtx.GetUserFee(ctx, repository.GetUserFeeParams{
		Userid:           creator.ID,
		Paymentchannelid: paymentChannel.ID,
	})
	if err != nil {
		s.log.Err(err).Msg("failed to get user fee")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 3. Calculate Amount & Fees
	 */

	donationAmount := req.Amount

	gatewayFee := (donationAmount*int64(fees.GatewayFeePercentageBps))/10000 +
		fees.GatewayFeeFixed

	platformFee := (donationAmount*int64(fees.PlatformFeePercentageBps))/10000 +
		fees.PlatformFeeFixed

	totalFeeFixed := fees.GatewayFeeFixed + fees.PlatformFeeFixed
	totalFeePercentageBps := fees.GatewayFeePercentageBps + fees.PlatformFeePercentageBps
	totalFeeAmount := gatewayFee + platformFee

	grossAmount := donationAmount + gatewayFee
	netAmount := grossAmount - totalFeeAmount

	/**
	 * 4. Create Donation Message
	 */

	donationMsgID, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate donation message id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	donationMsg, err := qtx.CreateDonationMessage(ctx,
		repository.CreateDonationMessageParams{
			ID:                donationMsgID,
			PayeeUserID:       req.PayeeUserID,
			PayerUserID:       utils.StringToPgTypeUUID(req.PayerUserID.String()),
			PayerName:         req.PayerName,
			Email:             req.Email,
			Message:           req.Message,
			MediaType:         repository.DonationMediaType(req.MediaType),
			MediaUrl:          utils.StringPtrToPgTypeText(req.MediaURL),
			MediaStartSeconds: utils.Int32PtrToPgTypeInt4(req.MediaStartSeconds),
			MaxPlaySeconds:    utils.MaxPlayedSeconds(int32(req.Amount), 500),
			PricePerSecond:    pgtype.Int8{Int64: 500, Valid: true},
			Amount:            donationAmount,
			Currency:          "IDR",
			Meta:              domain.JSONB{},
		},
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create donation message")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 5. Create Transaction Record
	 */

	transactionID, err := uuid.NewV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate transaction id")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	transactionResult, err := qtx.CreateTransaction(ctx,
		repository.CreateTransactionParams{
			ID:                       transactionID,
			DonationMessageID:        donationMsg.ID,
			PaymentChannelID:         paymentChannel.ID,
			PayeeUserID:              req.PayeeUserID,
			PayerUserID:              utils.StringToPgTypeUUID(req.PayerUserID.String()),
			GrossPaidAmount:          grossAmount,
			GatewayFeeFixed:          fees.GatewayFeeFixed,
			GatewayFeePercentageBps:  fees.GatewayFeePercentageBps,
			GatewayFeeAmount:         gatewayFee,
			PlatformFeeFixed:         fees.PlatformFeeFixed,
			PlatformFeePercentageBps: fees.PlatformFeePercentageBps,
			PlatformFeeAmount:        platformFee,
			FeeFixed:                 totalFeeFixed,
			FeePercentageBps:         totalFeePercentageBps,
			FeeAmount:                totalFeeAmount,
			NetAmount:                netAmount,
			Currency:                 "IDR",
			Status:                   repository.TransactionStatusINITIAL,
			ExternalReference:        pgtype.Text{Valid: false}, // di false dulu karena belum ada, akan di update setelah dapat dari midtrans
		},
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create transaction")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 7. Commit Transaction
	 */

	if err := tx.Commit(ctx); err != nil {
		s.log.Err(err).Msg("failed to commit transaction")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 8. Call Payment Gateway
	 */

	qrisResult, err := s.midtransClient.CreateQRIS(
		ctx,
		donationAmount,
		req.PayerName,
		req.Email,
		gatewayFee,
	)
	if err != nil {
		s.log.Err(err).Msg("failed to create qris midtrans")
		if err2 := s.q.UpdateTransactionStatus(ctx, repository.UpdateTransactionStatusParams{
			ID:     transactionResult.ID,
			Status: repository.TransactionStatusFAILED,
		}); err2 != nil {
			s.log.Err(err2).Msg("failed updating status after pg failure")
		}
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			"failed to create qris",
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 9. Update Transaction with External Reference
	 */

	if err := s.q.UpdateTransactionExternalReferenceAndStatus(ctx, repository.UpdateTransactionExternalReferenceAndStatusParams{
		ID:                transactionResult.ID,
		ExternalReference: utils.StringToPgTypeText(qrisResult.TransactionID),
		Status:            repository.TransactionStatusPENDING,
	}); err != nil {
		s.log.Err(err).Msg("failed to update transaction external reference")
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 10. Build Response
	 */

	return &CreateDonationResult{
		TransactionID: transactionID.String(),
		Amount:        transactionResult.GrossPaidAmount,
		Currency:      transactionResult.Currency,
		QrString:      qrisResult.QRString,
		Actions:       qrisResult.Actions,
	}, nil
}
