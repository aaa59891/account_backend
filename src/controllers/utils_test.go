package controllers_test

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/aaa59891/account_backend/src/db"
	"github.com/aaa59891/account_backend/src/inits"
	"github.com/aaa59891/account_backend/src/routers"

	"github.com/gin-gonic/gin"
)

const (
	urlPrefix           = ""
	apiToken            = "api_token"
	errStrUnmarshalBody = "could not unmarshal response body: %v"
	errStrCreateRequest = "could not create request: %v"
	errStrMarshal       = "could not marshal account form: %v"
	errStrQueryDb       = "could not query db: %v"
	errStrCreateData    = "could not create data: %v"
)

var router *gin.Engine
var token string

type errResponse struct {
	Message string
}

type responseBody struct {
	Data interface{}
}

func init() {
	gin.DisableConsoleColor()
	gin.SetMode(gin.TestMode)

	router = gin.New()
	router.Use(GinRecover)

	routers.SetRoutes(router)
	inits.CreateTable()
}

func CreateModel(model interface{}, t *testing.T) {
	if err := db.DB.Create(model).Error; err != nil {
		t.Fatalf(errStrCreateData, err)
	}
}

func DeleteModel(model, id interface{}, pkName string, t *testing.T) {
	if err := db.DB.Unscoped().Delete(model, pkName+" = ?", id).Error; err != nil {
		t.Fatalf("could not delete model: %v", err)
	}
}

func GinRecover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Next()
		}
	}()
	c.Next()
}

func GetJsonBody(obj interface{}) io.Reader {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(b)
}

func SetUploadFileBody(writer *multipart.Writer, paramName, path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return err
	}

	part.Write(fileContents)

	return nil
}

func GetFormBody(obj interface{}) *bytes.Buffer {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	defer writer.Close()
	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < v.NumField(); i++ {
		kind := v.Field(i).Kind()
		tag := v.Type().Field(i).Tag
		formTag := tag.Get("form")
		value := v.Field(i)

		if len(formTag) == 0 {
			continue
		}

		switch kind {
		case reflect.Struct:
		default:
			writer.WriteField(formTag, fmt.Sprint(value))
		}
	}
	return body
}

type testHelper struct {
	Response        *httptest.ResponseRecorder
	Request         *http.Request
	ResponseErrBody errResponse
	Errs            []error
	T               *testing.T
}

func (th *testHelper) CheckErrors() *testHelper {
	if len(th.Errs) > 0 {
		th.T.Errorf("send request had errors: %v", th.Errs)
	}
	return th
}

func (th *testHelper) AddError(err error) *testHelper {
	if th.Errs == nil {
		th.Errs = make([]error, 0)
	}
	th.Errs = append(th.Errs, err)
	return th
}

func (th *testHelper) SetRequest(method, url string, body io.Reader) *testHelper {
	defer th.CheckErrors()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return th.AddError(err)
	}
	th.Request = req
	return th
}

func (th *testHelper) SendRequest(router http.Handler) *testHelper {
	defer th.CheckErrors()
	if th.Response == nil {
		th.Response = httptest.NewRecorder()
	}
	if th.Request == nil {
		th.AddError(errors.New("request is nil"))
		return th
	}

	router.ServeHTTP(th.Response, th.Request)
	return th
}

func (th *testHelper) SetHeader(key, value string) *testHelper {
	defer th.CheckErrors()
	if th.Request == nil {
		return th.AddError(errors.New("request is nil"))
	}
	th.Request.Header.Set(key, value)
	return th
}

func (th *testHelper) DecodeResponseBody(body interface{}) *testHelper {
	defer th.CheckErrors()
	if th.Response == nil {
		th.AddError(errors.New("could not decode response body: response is nil"))
		return th
	}
	if err := json.Unmarshal([]byte(th.Response.Body.String()), body); err != nil {
		th.AddError(fmt.Errorf(errStrUnmarshalBody, err))
	}
	return th
}

func (th *testHelper) DecodeErrResponseBody() *testHelper {
	defer th.CheckErrors()
	if err := json.Unmarshal([]byte(th.Response.Body.String()), &th.ResponseErrBody); err != nil {
		th.AddError(fmt.Errorf(errStrUnmarshalBody, err))
	}
	return th
}

func GetTestHelper(t *testing.T) *testHelper {
	return &testHelper{T: t}
}

func GetQueryString(obj interface{}) string {
	queryStrings := make([]string, 0)
	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < v.NumField(); i++ {
		kind := v.Field(i).Kind()
		tag := v.Type().Field(i).Tag
		formTag := tag.Get("form")
		value := v.Field(i)

		if len(formTag) == 0 {
			continue
		}

		switch kind {
		case reflect.Struct:
			t, ok := value.Interface().(time.Time)
			dateFormat := tag.Get("time_format")
			if !ok || len(dateFormat) == 0 || t.IsZero() {
				continue
			}
			queryStrings = append(queryStrings, fmt.Sprint(formTag, "=", t.Format(dateFormat)))
		default:
			if value.Interface() == reflect.Zero(value.Type()).Interface() {
				continue
			}
			queryStrings = append(queryStrings, fmt.Sprint(formTag, "=", value))
		}
	}

	return strings.Join(queryStrings, "&")
}

type testModel struct {
	name   string
	status int
	err    string
	model  interface{}
}
