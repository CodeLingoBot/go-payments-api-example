package testutil

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

// NewUUID creates a random V4 UUID and abort the test if there is an error.
func NewUUID(t *testing.T) uuid.UUID {
	var u, err = uuid.NewV4()
	require.NoError(t, err)

	return u
}
