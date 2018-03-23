package controllers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

var defaultRecord = models.Record{Email: "test@test.com", CategoryId: 1, Amount: 200, Date: time.Now()}

func TestCreateRecord(t *testing.T) {
	tm := []testModel{
		{
			"create a new record",
			http.StatusOK,
			"",
			defaultRecord,
		},
		{
			"create a new record with empty email",
			http.StatusBadRequest,
			models.ErrNoEmail.Error(),
			models.Record{CategoryId: 1, Amount: 200},
		},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			record := tc.model.(models.Record)

			buf := GetJsonBody(record)

			th := GetTestHelper(tt).SetRequest(http.MethodPost, urlPrefix+"/record", buf).SendRequest(router)

			assert.Equal(tt, th.Response.Code, tc.status)

			if len(tc.err) > 0 {
				th.DecodeErrResponseBody()
				assert.Equal(tt, th.ResponseErrBody.Message, tc.err)
				return
			}

			dbRecord := models.Record{}

			if err := db.DB.First(&dbRecord, "email = ? and amount = ? and category_id = ?", record.Email, record.Amount, record.CategoryId).Error; err != nil {
				tt.Fatalf("could not find the record: %v", err)
			}

			DeleteModel(&models.Record{}, dbRecord.Id, "id", tt)
		})

	}
}

func TestFetchRecords(t *testing.T) {
	data := []models.Record{
		{Email: defaultRecord.Email, Date: defaultRecord.Date, Amount: 100},
		{Email: defaultRecord.Email, Date: defaultRecord.Date, Amount: 200},
		{Email: defaultRecord.Email, Date: defaultRecord.Date, Amount: 300},
	}

	if err := models.Transactional(func(tx *gorm.DB) error {
		for _, d := range data {
			if err := tx.Create(&d).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("could not create record: %v", err)
	}

	tm := []struct {
		testModel
		length int
	}{
		{
			testModel{
				"fetch record in date mode",
				http.StatusOK,
				"",
				models.RecordForm{Email: defaultRecord.Email, Date: defaultRecord.Date, Mode: models.RECORD_FETCH_MODE_DATE},
			},
			len(data),
		},
		{
			testModel{
				"fetch record in month mode",
				http.StatusOK,
				"",
				models.RecordForm{Email: defaultRecord.Email, Date: defaultRecord.Date, Mode: models.RECORD_FETCH_MODE_MONTH},
			},
			len(data),
		},
		{
			testModel{
				"fetch record with different month in month mode",
				http.StatusOK,
				"",
				models.RecordForm{Email: defaultRecord.Email, Date: defaultRecord.Date.Add(time.Hour * 24 * 31), Mode: models.RECORD_FETCH_MODE_MONTH},
			},
			0,
		},
		{
			testModel{
				"fetch record in all mode",
				http.StatusOK,
				"",
				models.RecordForm{Email: defaultRecord.Email, Date: time.Time{}, Mode: models.RECORD_FETCH_MODE_ALL},
			},
			len(data),
		},
	}

	for _, tc := range tm {
		t.Run(tc.name, func(tt *testing.T) {
			rf := tc.model.(models.RecordForm)

			th := GetTestHelper(tt).SetRequest(http.MethodGet, urlPrefix+fmt.Sprintf("/record?email=%s&date=%s&mode=%s", rf.Email, rf.Date.Format("2006-01-02"), rf.Mode), nil).SendRequest(router)

			assert.Equal(tt, th.Response.Code, tc.status)

			body := struct {
				Data []models.Record
			}{}
			th.DecodeResponseBody(&body)
			assert.Equal(tt, len(body.Data), tc.length)
		})
	}

	if err := models.Transactional(func(tx *gorm.DB) error {
		for _, d := range data {
			if err := tx.Delete(&d).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("could not delete record: %v", err)
	}
}
