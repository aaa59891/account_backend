package controllers_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

func TestSingUp(t *testing.T) {
	existUser := models.User{Email: "exist@test.com", Password: "justtestpassword"}
	if err := db.DB.Create(&existUser).Error; err != nil {
		t.Fatalf(errStrCreateData, err)
	}

	tm := []testModel{
		{"sign up a new account", http.StatusOK, "", models.User{Email: "chong@test.com", Password: "testPassword"}},
		{"sign up an exist account", http.StatusBadRequest, models.ErrEmailExist.Error(), existUser},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			newUser := tc.model.(models.User)

			buf := GetJsonBody(newUser)

			th := GetTestHelper(tt).SetRequest(http.MethodPost, urlPrefix+"/user", buf).SendRequest(router)
			assert.Equal(tt, th.Response.Code, tc.status)

			if len(tc.err) > 0 {
				th.DecodeErrResponseBody()
				assert.Equal(t, th.ResponseErrBody.Message, tc.err)
				return
			}

			dbUser := models.User{}
			if err := db.DB.First(&dbUser, "email = ?", strings.ToLower(newUser.Email)).Error; err != nil {
				tt.Fatalf("could not find the new user: %v", err)
			}
			DeleteModel(&models.User{}, dbUser.Id, "id", tt)
		})
	}

	DeleteModel(&existUser, existUser.Id, "id", t)
}

func TestSingIn(t *testing.T) {

	user := models.User{Email: "chong@test.com", Password: "testPassword"}

	signUp := GetJsonBody(user)

	signUpTh := GetTestHelper(t).SetRequest(http.MethodPost, urlPrefix+"/user", signUp).SendRequest(router)

	if len(signUpTh.Errs) > 0 {
		t.Fatalf("could not sign up: %v", signUpTh.Errs)
	}
	tm := []testModel{
		{"test sign in", http.StatusOK, "", user},
		{"test sign in with wrong password", http.StatusBadRequest, models.ErrWrongPassword.Error(), models.User{Email: user.Email, Password: "wrongPassword"}},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			u := tc.model.(models.User)
			signInBuf := GetJsonBody(u)
			signInTh := GetTestHelper(tt).SetRequest(http.MethodPost, urlPrefix+"/signin", signInBuf).SendRequest(router)

			assert.Equal(tt, signInTh.Response.Code, tc.status)

			if len(tc.err) > 0 {
				signInTh.DecodeErrResponseBody()
				assert.Equal(tt, signInTh.ResponseErrBody.Message, tc.err)
				return
			}
		})
	}
	DeleteModel(&models.User{}, user.Email, "email", t)
}
