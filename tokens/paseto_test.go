package tokens

import (
	"testing"
	"time"

	"github.com/sirjager/trueauth/utils"
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
	data := PayloadData{UserID: utils.RandomUUID(), UserEmail: utils.RandomEmail()}

	hash, payload, err := builder.CreateToken(data, time.Second*10)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.Equal(t, data.UserEmail, payload.Payload.UserEmail)

	// Now verify
	payload, err = builder.VerifyToken(hash)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.NotEmpty(t, payload.Payload.UserID)
	require.NotEmpty(t, payload.Payload.UserEmail)
	require.Equal(t, data.UserID, payload.Payload.UserID)
	require.Equal(t, data.UserEmail, payload.Payload.UserEmail)

	// Verify Token
	// with expired token
	hash, payload, err = builder.CreateToken(data, time.Microsecond)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.Id)
	require.NotEmpty(t, payload.IssuedAt)
	require.NotEmpty(t, payload.ExpiresAt)
	require.Equal(t, data.UserID, payload.Payload.UserID)
	require.Equal(t, data.UserEmail, payload.Payload.UserEmail)
	// Now verify with expired token
	payload, err = builder.VerifyToken(hash)
	require.Error(t, err)
	require.Empty(t, payload)
	require.EqualError(t, ErrExpiredToken, err.Error())
}
