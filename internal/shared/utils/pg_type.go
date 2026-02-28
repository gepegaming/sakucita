package utils

import (
	"encoding/json"
	"errors"
	"sakucita/internal/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgTypeText(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: str, Valid: true}
}

func StringPtrToPgTypeText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
func StringPtrToPgTypeUUID(str *string) pgtype.UUID {
	if str == nil || *str == "" || *str == "00000000-0000-0000-0000-000000000000" {
		return pgtype.UUID{Valid: false}
	}
	u, err := uuid.Parse(*str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func StringToPgTypeUUID(str string) pgtype.UUID {
	if str == "" || str == "00000000-0000-0000-0000-000000000000" {
		return pgtype.UUID{Valid: false}
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func UUIDPtrToPgTypeUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil || *u == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func Int32ToPgTypeInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func Int32PtrToPgTypeInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}

func Int64ToPgTypeInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

func Int64PtrToPgTypeInt8(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

// pgtype
func PgTypeUUIDToUUIDPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}

	id := uuid.UUID(u.Bytes)
	return &id
}

func PgTypeUUIDToUUID(u pgtype.UUID) (uuid.UUID, error) {
	if !u.Valid {
		return uuid.UUID{}, domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid uuid",
			errors.New("uuid is null or undefined"),
		)
	}

	return uuid.UUID(u.Bytes), nil
}

func PgTypeTextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func PgTypeTextToString(t pgtype.Text) (string, error) {
	if !t.Valid {
		return "", domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid text",
			errors.New("text is null or undefined"),
		)
	}
	return t.String, nil
}

func PgTypeInt4ToInt32(i pgtype.Int4) (*int32, error) {
	if !i.Valid {
		return nil, domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid int4",
			errors.New("int4 is null or undefined"),
		)
	}
	return &i.Int32, nil
}

func PgTypeInt4ToInt32Ptr(i pgtype.Int4) (*int32, error) {
	if !i.Valid {
		return nil, domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid int4",
			errors.New("int4 is null or undefined"),
		)
	}
	return &i.Int32, nil
}

func PgTypeInt8ToInt64Ptr(i pgtype.Int8) (*int64, error) {
	if !i.Valid {
		return nil, domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid int8",
			errors.New("int8 is null or undefined"),
		)
	}
	return &i.Int64, nil
}

func PgTypeInt8ToInt64(i pgtype.Int8) (int64, error) {
	if !i.Valid {
		return 0, domain.NewAppError(
			fiber.StatusBadRequest,
			"invalid int8",
			errors.New("int8 is null or undefined"),
		)
	}
	return i.Int64, nil
}

// jsonb
func JSONBytesToMap(b []byte) (map[string]any, error) {
	if len(b) == 0 {
		return nil, nil
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}
