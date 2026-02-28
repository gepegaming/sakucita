package service

import (
	"context"

	"sakucita/internal/domain"
	"sakucita/internal/infra/midtrans"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"

	feeService "sakucita/internal/app/fee/service"
	PaymentChannelService "sakucita/internal/app/paymentChannel/service"
	transactionService "sakucita/internal/app/transaction/service"
	userService "sakucita/internal/app/user/service"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	DefaultPricePerSecond int64 = 500
	DefaultCurrency             = "IDR"
)

type service struct {
	db                    *pgxpool.Pool
	q                     *repository.Queries
	log                   zerolog.Logger
	midtransClient        midtrans.MidtransClient
	trasanctionService    transactionService.TransactionService
	feeService            feeService.FeeService
	userService           userService.UserService
	paymentChannelService PaymentChannelService.PaymentChannelService
}

type DonationService interface {
	CreateDonation(ctx context.Context, req CreateDonationCommand) (CreateDonationResult, error)
	createDonationMessageWithTx(ctx context.Context, qtx repository.Querier, cmd CreateDonationMessageCommand) (domain.DonationMessage, error)
}

func NewService(
	db *pgxpool.Pool,
	q *repository.Queries,
	log zerolog.Logger,
	midtransClient midtrans.MidtransClient,
	transactionService transactionService.TransactionService,
	feeService feeService.FeeService,
	userService userService.UserService,
	paymentChannelService PaymentChannelService.PaymentChannelService,
) DonationService {
	return &service{db, q, log, midtransClient, transactionService, feeService, userService, paymentChannelService}
}

func (s *service) CreateDonation(
	ctx context.Context,
	req CreateDonationCommand,
) (res CreateDonationResult, err error) {

	/**
	 * 1. Setup Database Transaction
	 */
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.log.Err(err).Msg("failed to begin transaction")
		return res, domain.NewAppError(
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

	creator, err := s.userService.GetByIDWithTx(ctx, qtx, req.PayeeUserID)
	if err != nil {
		return res, err
	}

	paymentChannel, err := s.paymentChannelService.GetPaymentChannelByCodeWithTx(ctx, qtx, req.PaymentChannel)
	if err != nil {
		return res, err
	}

	fees, err := s.feeService.GetUserFeesWithTx(ctx, qtx, creator.ID, paymentChannel.ID)
	if err != nil {
		return res, err
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

	donationMsg, err := s.createDonationMessageWithTx(ctx, qtx, CreateDonationMessageCommand{
		PayeeUserID:       creator.ID,
		PayerUserID:       req.PayerUserID,
		PayerName:         req.PayerName,
		Email:             req.Email,
		Message:           req.Message,
		MediaType:         repository.DonationMediaType(req.MediaType),
		MediaUrl:          req.MediaURL,
		MediaStartSeconds: req.MediaStartSeconds,
		PricePerSecond:    DefaultPricePerSecond,
		GrossPaidAmount:   grossAmount,
		Amount:            donationAmount,
		Currency:          DefaultCurrency,
		Meta:              []byte("{}"),
	})
	if err != nil {
		return res, err
	}

	/**
	 * 5. Create Transaction Record
	 */

	transactionResult, err := s.trasanctionService.CreateWithTx(ctx, qtx,
		transactionService.CreateTransactionCommand{
			DonationMessageID:        donationMsg.ID,
			PaymentChannelID:         paymentChannel.ID,
			PayeeUserID:              req.PayeeUserID,
			PayerUserID:              req.PayerUserID,
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
			Currency:                 DefaultCurrency,
			Status:                   domain.TransactionStatusInitial,
			ExternalReference:        pgtype.Text{Valid: false}, // di false dulu karena belum ada, akan di update setelah dapat dari midtrans
		},
	)
	if err != nil {
		return res, err
	}

	/**
	 * 7. Commit Transaction
	 */

	if err := tx.Commit(ctx); err != nil {
		s.log.Err(err).Msg("failed to commit transaction for create donation")
		return res, domain.NewAppError(
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
		if err2 := s.trasanctionService.MarkAsFailed(ctx, transactionResult.ID); err2 != nil {
			s.log.Err(err2).Msg("failed updating status after pg failure")
		}
		return res, domain.NewAppError(
			fiber.StatusInternalServerError,
			"failed to create qris",
			domain.ErrInternalServerError,
		)
	}

	/**
	 * 9. Update Transaction with External Reference
	 */

	if err := s.trasanctionService.UpdateExternalReferenceAndStatus(ctx, transactionResult.ID, qrisResult.TransactionID, domain.TransactionStatusPending); err != nil {
		return res, err
	}

	/**
	 * 10. Build Response
	 */

	return CreateDonationResult{
		TransactionID: transactionResult.ID.String(),
		Amount:        transactionResult.GrossPaidAmount,
		Currency:      transactionResult.Currency,
		QrString:      qrisResult.QRString,
		Actions:       qrisResult.Actions,
	}, nil
}

func (s *service) createDonationMessageWithTx(ctx context.Context, qtx repository.Querier, cmd CreateDonationMessageCommand) (domain.DonationMessage, error) {
	donationMsgID, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate donation message id")
		return domain.DonationMessage{}, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	donationMsg, err := qtx.CreateDonationMessage(ctx,
		repository.CreateDonationMessageParams{
			ID:                donationMsgID,
			PayeeUserID:       cmd.PayeeUserID,
			PayerUserID:       utils.UUIDPtrToPgTypeUUID(cmd.PayerUserID),
			PayerName:         cmd.PayerName,
			Email:             cmd.Email,
			Message:           cmd.Message,
			MediaType:         repository.DonationMediaType(cmd.MediaType),
			MediaUrl:          utils.StringPtrToPgTypeText(cmd.MediaUrl),
			MediaStartSeconds: utils.Int32PtrToPgTypeInt4(cmd.MediaStartSeconds),
			MaxPlaySeconds:    maxPlayedSeconds(cmd.Amount, DefaultPricePerSecond),
			PricePerSecond:    pgtype.Int8{Int64: DefaultPricePerSecond, Valid: true},
			Amount:            cmd.GrossPaidAmount,
			Currency:          DefaultCurrency,
			Meta:              cmd.Meta,
		},
	)

	if err != nil {
		s.log.Err(err).Msg("failed to create donation message")
		return domain.DonationMessage{}, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	return donationMessageRepoToDomain(donationMsg), nil
}

// utils function to calculate max played seconds based on amount and price per second
func maxPlayedSeconds(amount, pricePerSecond int64) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(amount / pricePerSecond), Valid: true}
}
