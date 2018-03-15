package controllers_test

import (
	"net/http"
	"testing"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

func TestCreateCategory(t *testing.T) {
	tm := []struct {
		testModel
		checkDatabase bool
	}{
		{
			testModel{
				"create category",
				http.StatusOK,
				"",
				models.Category{Email: "chong@email.com", Name: "category name"},
			},
			true,
		},
		{
			testModel{
				"create category with empty email",
				http.StatusBadRequest,
				"",
				models.Category{Email: "", Name: "category name"},
			},
			false,
		},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			category := tc.model.(models.Category)

			buf := GetJsonBody(category)

			th := GetTestHelper(tt).SetRequest(http.MethodPost, urlPrefix+"/category", buf).SendRequest(router)

			assert.Equal(tt, th.Response.Code, tc.status)
			if !tc.checkDatabase {
				return
			}
			dbCategory := models.Category{}

			if err := db.DB.First(&dbCategory, "email = ? and name = ?", category.Email, category.Name).Error; err != nil {
				tt.Fatalf("could not find the category: %v", err)
			}

			DeleteModel(&models.Category{}, dbCategory.Id, "id", tt)
		})
	}
}
