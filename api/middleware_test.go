package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/token"
)

// create sample authorization token assign into request header
func addAuthorization(t *testing.T, request *http.Request, tkMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	token, err := tkMaker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Add(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddle(t *testing.T) {
	testcases := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tkMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tkMaker token.Maker) {
				addAuthorization(t, request, tkMaker, authorizationTypeBearer, "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorize",
			setUpAuth: func(t *testing.T, request *http.Request, tkMaker token.Maker) {
				addAuthorization(t, request, tkMaker, "unauthorized", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			// make a mock http request (include the authorization header)
			server := newTestServer(t, nil)

			authPath := "/auth"
			// create mock server with test method & handler of function
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			tc.setUpAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)

		})
	}
}
