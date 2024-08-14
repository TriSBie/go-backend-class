package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_sqlc "simple_bank.sqlc.dev/app/db/mock"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/util"
)

/*
gomock is mocking framework with go, integrates with testing package
*/

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Assert that Bar() is invoked.

	store := mock_sqlc.NewMockStore(ctrl)
	// Expect the getAccountById (any context, should be equal with the account generated by random)
	// GetAccountById(ctx context, id int64) -> (Account acc, err Error)
	store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	// test server and send request
	server := NewServer(store)
	// recorder will record the response result from serveHTTP
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	log.Println("url : ", url)
	// define new request include method and url to request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	log.Println("err", err)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, req)

	log.Println("Res body", recorder.Body)
	log.Printf("Res body type %T\n", recorder.Body)
	// check response
	// require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func TestGetAccountAPI_Refactoring(t *testing.T) {
	account := randomAccount()

	/*
	  # create a build_stubs with array of anonymous struct
	  @name: using for test the response code
	  @accountID: parameter required to test
	  @buildStubs: define a function for pre-configuration setting the test case
	  @checkResponse: validate or inspect the results response from http-test
	*/
	test_cases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mock_sqlc.MockStore) {
				// Expect the getAccountById (any context, should be equal with the account generated by random)
				// GetAccountById(ctx context, id int64) -> (Account acc, err Error)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: 0,
			buildStubs: func(store *mock_sqlc.MockStore) {
				// Expect the getAccountById (any context, should be equal with the account generated by random)
				// GetAccountById(ctx context, id int64) -> (Account acc, err Error)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(0)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// requireBodyMatchAccount(t, recorder.Body)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mock_sqlc.MockStore) {
				// Expect the getAccountById (any context, should be equal with the account generated by random)
				// GetAccountById(ctx context, id int64) -> (Account acc, err Error)
				store.EXPECT().GetAccountById(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// requireBodyMatchAccount(t, recorder.Body)
			},
		},
	}

	for i := range test_cases {
		tc := test_cases[i]
		// run subset of test case with given name and executes parallel testing using with goroutine
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Assert that Bar() is invoked.

			store := mock_sqlc.NewMockStore(ctrl)
			tc.buildStubs(store)

			// test server and send request
			server := NewServer(store)
			// recorder will record the response result from serveHTTP
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			fmt.Printf("Constructed URL: %s\n", url)
			// define new request include method and url to request, body
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			fmt.Printf("Redirect Location: %s\n", recorder.Header().Get("Location"))
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(316, 323),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var goAccount db.Account
	// parse from json into db account type
	err = json.Unmarshal(data, &goAccount)
	require.NoError(t, err)
	require.Equal(t, account, goAccount)

}
