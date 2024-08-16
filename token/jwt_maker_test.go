package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
)

func TestJWTMaker(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration) // 1 minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)

	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpired.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// create sample payload
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	// create sample token sing UnSageAllowNoneSignature
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// create jwt maker with secret key
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	// compare non-sense jwt token generated manually compare to difference signature
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
