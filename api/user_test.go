package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "simple_bank.sqlc.dev/app/db/mock"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/token"
	"simple_bank.sqlc.dev/app/util"
)

// create custom matcher func

type argMatcherCreateUserParam struct {
	arg      db.CreateUserParams
	password string
}

func (e argMatcherCreateUserParam) Matches(x interface{}) bool {
	// using type assertion and check error
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// checking the password & hashPassword
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	// assign hashed password into struct with field arg
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return argMatcherCreateUserParam{
		arg:      arg,
		password: password,
	}
}

func (e argMatcherCreateUserParam) String() string {
	return fmt.Sprintf("is equal to %v (%T)", e.password, e.arg.HashedPassword)
}

func TestCreateUserAPi(t *testing.T) {
	user, password := randomUser(t)

	test_cases := []struct {
		name          string
		body          gin.H //shorthand for map[string]any,
		setUpAuth     func(t *testing.T, request *http.Request, tkMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			setUpAuth: func(t *testing.T, request *http.Request, tkMaker token.Maker) {
				addAuthorization(t, request, tkMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			// mockStore is generated from structure Queries methods
			buildStubs: func(store *mockdb.MockStore) {
				// change gomock params matcher more strictly
				arg := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
				}
				fmt.Println("Password before: ", user.HashedPassword)
				// EXPECT ANY VALUE - expect function CreateUser called -> return user with nil err
				// Compare CreateUser(context, arg) -> args does matches with the actual execution or not
				// passing an argument of createUserParams with the initial password

				// WHY USING CUSTOM goMock equal custom checker ? -> Since createUser function handler calls every new hashPassword as each time executed.
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// validate the response code
				require.Equal(t, http.StatusCreated, recorder.Code)
				fmt.Println("Password after: ", recorder.Body)
				requireBodyMatcherUser(t, recorder.Body, user)
			},
		},
	}

	for i := range test_cases {
		// doing required steps before testing
		tc := test_cases[i]

		t.Run(tc.name, func(t *testing.T) {

			// init gomock with controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Assert that ctrl.Finished is invoked.

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			// init recorder
			recorder := httptest.NewRecorder()

			urlPath := "/users"

			// convert from map[string]any types -> bytes value and passing into reader header
			data, err := json.Marshal(tc.body)

			require.NoError(t, err)

			// making new HTTP request for making dump test with http
			req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setUpAuth(t, req, server.tokenMaker)
			// passing recorder var to record results response from the serveHTTP making
			// [*] Invoke a handler of method to being test with stubs has already declared above
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)

		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)

	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	username := util.RandomString(6)
	userCreate := db.User{
		Username:       username,
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          username + "@gmail.com",
	}
	require.NotEmpty(t, userCreate)
	require.Equal(t, username, userCreate.Username)

	return userCreate, password
}

func requireBodyMatcherUser(t *testing.T, body *bytes.Buffer, user db.User) {
	// convert into bytes
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// convert bytes into JSON using UnMmarshall -> from bytes into a specific data types
	var userArg db.User
	err = json.Unmarshal(data, &userArg)
	require.NoError(t, err)

	userArgRes := newUserResponse(userArg)
	userExpectedRes := newUserResponse(user)

	require.Equal(t, userArgRes, userExpectedRes)
}
