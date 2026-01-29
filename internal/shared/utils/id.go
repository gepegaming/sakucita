package utils

import "github.com/google/uuid"

func GenerateUUIDV7() (uuid.UUID, error) {
	return uuid.NewV7()
}

func GenerateUUIDV7MustNew() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
