package controllers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

var defaultCategory = models.Category{Email: "chong@email.com", Name: "category name"}

func TestCreateCategory(t *testing.T) {
	tm := []testModel{
		{
			"create category",
			http.StatusOK,
			"",
			defaultCategory,
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
	category := defaultCategory
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

func TestDeleteCategory(t *testing.T) {
	category := defaultCategory

	if err := db.DB.Create(&category).Error; err != nil {
		t.Fatalf("could not create category: %v", err)
	}

	tm := []testModel{
		{"delete a category", http.StatusOK, "", category},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			c := tc.model.(models.Category)

			th := GetTestHelper(tt).SetRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%d", urlPrefix, "category", c.Id), nil).SendRequest(router)
			assert.Equal(tt, th.Response.Code, tc.status)

			dbCategory := models.Category{}
			err := db.DB.Find(&dbCategory, "id = ?", c.Id).Error
			assert.Equal(tt, err, gorm.ErrRecordNotFound)
		})
	}
}

func TestFetchCategories(t *testing.T) {
	categories := []models.Category{
		{Email: defaultCategory.Email, Name: "test1"},
		{Email: defaultCategory.Email, Name: "test2"},
		{Email: defaultCategory.Email, Name: "test3"},
	}
	if err := models.Transactional(func(tx *gorm.DB) error {
		for _, category := range categories {
			if err := tx.Create(&category).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("could not create category: %v", err)
	}

	tm := []testModel{
		{
			"fetch categories",
			http.StatusOK,
			"",
			defaultCategory.Email,
		},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			email := tc.model.(string)
			th := GetTestHelper(tt).SetRequest(http.MethodGet, urlPrefix+"/category/"+email, nil).SendRequest(router)

			assert.Equal(tt, th.Response.Code, tc.status)

			if len(tc.err) > 0 {
				th.DecodeErrResponseBody()
				assert.Equal(tt, th.ResponseErrBody.Message, tc.err)
				return
			}
			body := struct {
				Data []models.Category
			}{}
			th.DecodeResponseBody(&body)
			assert.Equal(tt, len(body.Data), len(categories))
		})
	}
	if err := models.Transactional(func(tx *gorm.DB) error {
		for _, category := range categories {
			if err := tx.Delete(category).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("could not delete category: %v", err)
	}

}
