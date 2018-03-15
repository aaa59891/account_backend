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
