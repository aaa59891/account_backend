package controllers_test

import (
	"net/http"
	"testing"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

func TestCreateCategory(t *testing.T) {
	tm := []testModel{
		{
			"create category",
			http.StatusOK,
			"",
			models.Category{Email: "chong@email.com", Name: "category name"},
		},
		{
			"create category with empty email",
			http.StatusBadRequest,
			models.ErrNoEmail.Error(),
			models.Category{Email: "", Name: "category name"},
		},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			category := tc.model.(models.Category)

			buf := GetJsonBody(category)

			th := GetTestHelper(tt).SetRequest(http.MethodPost, urlPrefix+"/category", buf).SendRequest(router)

			assert.Equal(tt, th.Response.Code, tc.status)
			if len(tc.err) > 0 {
				th.DecodeErrResponseBody()
				assert.Equal(tt, th.ResponseErrBody.Message, tc.err)
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

func TestUpdateCategory(t *testing.T) {
	category := models.Category{Email: "test@test.com", Name: "testname"}
	if err := db.DB.Create(&category).Error; err != nil {
		t.Fatalf("could not create category: %v", err)
	}
	tm := []testModel{
		{"update category", http.StatusOK, "", models.Category{Id: category.Id, Name: "newname"}},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			c := tc.model.(models.Category)

			buf := GetJsonBody(c)

			th := GetTestHelper(tt).SetRequest(http.MethodPut, urlPrefix+"/category", buf).SendRequest(router)
			assert.Equal(tt, th.Response.Code, tc.status)

			dbCategory := models.Category{}
			if err := db.DB.Find(&dbCategory, "id = ?", c.Id).Error; err != nil {
				tt.Fatalf("could not find the category: %v", err)
			}

			assert.Equal(tt, dbCategory.Name, c.Name)
			assert.Equal(tt, dbCategory.Email, category.Email)
		})
	}

	DeleteModel(&models.Category{}, category.Id, "id", t)
}
