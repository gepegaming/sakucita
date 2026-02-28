package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func Int32Ptr(v int) *int32 {
	i := int32(v)
	return &i
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func ParseRupiahAmount(s string) (int64, error) {
	parts := strings.Split(s, ".")
	return strconv.ParseInt(parts[0], 10, 64)
}

func MapToByteSlice(m map[string]any) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map: %w", err)
	}
	return b, nil
}
