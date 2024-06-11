package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID string

// New generates a new random UUID string.
// The UUID string generated using this method conform to
// version 4
func NewV4() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("unable to uuid.New() : %w", err))
	}

	return id.String()
}

// New generates a new random UUID string.
// The UUID string generated using this method conform to
// version 7
func NewV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("unable to uuid.New() : %w", err))
	}

	return id.String()
}
