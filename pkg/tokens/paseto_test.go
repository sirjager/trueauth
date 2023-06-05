package tokens

import (
	"testing"
	"time"

	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoBuilder(t *testing.T) {
	builder, err := NewPasetoBuilder(small_secret_key)
	require.Error(t, err)
	require.Empty(t, builder)

	builder, err = NewPasetoBuilder(valid_secret_key)
	require.NoError(t, err)
	require.NotEmpty(t, builder)

	// Create Token
	data := PayloadData{UserID: utils.UUID_XID(), UserEmail: utils.RandomEmail()}

	hash, payload, err := builder.CreateToken(data, time.Second*10)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.ID)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.Equal(t, data.UserEmail, payload.Data.UserEmail)

	// Now verify
	payload, err = builder.VerifyToken(hash)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.ID)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.NotEmpty(t, payload.Data.UserID)
	require.NotEmpty(t, payload.Data.UserEmail)
	require.Equal(t, data.UserID, payload.Data.UserID)
	require.Equal(t, data.UserEmail, payload.Data.UserEmail)

	// Verify Token
	// with expired token
	hash, payload, err = builder.CreateToken(data, time.Microsecond)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.ID)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.Equal(t, data.UserID, payload.Data.UserID)
	require.Equal(t, data.UserEmail, payload.Data.UserEmail)
	// Now verify with expired token
	payload, err = builder.VerifyToken(hash)
	require.Error(t, err)
	require.Empty(t, payload)
	require.EqualError(t, ErrExpiredToken, err.Error())
}
