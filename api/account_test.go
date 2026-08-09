package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"  // fix: was missing gomock subpackage
	mockdb "github.com/simpleBank/db/mock"
	db "github.com/simpleBank/db/sqlc"
	"github.com/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)  // fix: NewContoller → NewController
	defer ctrl.Finish()              // fix: defer.ctrl → defer ctrl

	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)  // fix: fmtSprintf → fmt.Sprintf
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)  // fix: ServerHttp → ServeHTTP

	// check response
	require.Equal(t, http.StatusOK, recorder.Code)  // fix: StatusOk → StatusOK
}

func randomAccount() db.Account {
	return db.Account{              // fix: missing return keyword
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}