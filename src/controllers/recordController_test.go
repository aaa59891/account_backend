package controllers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/aaa59891/account_backend/src/db"

	"github.com/magiconair/properties/assert"

	"github.com/aaa59891/account_backend/src/models"
)

var defaultRecord = models.Record{Email: "test@tets.com", CategoryId: 1, Amount: 200, Date: time.Now()}

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
