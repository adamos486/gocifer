package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cobb.com/server/database"
	. "cobb.com/server/events/controllers"
	fakeDb "cobb.com/server/fakes/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
	"cobb.com/server/events/services"
	"bytes"
)

func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
	}
	return r
}

func testHTTPResponse(r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	//Create a new response recorder
	w := httptest.NewRecorder()
	//Create the service and process the above request.
	r.ServeHTTP(w, req)

	recorded := f(w)
	Expect(recorded).To(Equal(true))
}

type AddEventApiResponse struct {
	Code    int         `json:"code"`
	Created eventObject `json:"created"`
}

type GetEventsApiResponse struct {
	Code int           `json:"code"`
	List []eventObject `json:"list"`
}

type eventObject struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:",string"`
}

func (add *AddEventApiResponse) doesItPassTheContract() bool {
	return add.Code == 200 && add.Created.doesItPassTheContract()
}

func (cr *eventObject) doesItPassTheContract() bool {
	return cr.Id != -1 && cr.Name != "" && cr.Description != "" && cr.Date != time.Time{}
}

var _ = Describe("Events Controller", func() {
	var (
		client          database.Client
		fakeDbClient    *fakeDb.FakeDBClient
		service         *services.EventsServiceClient
		db              *sql.DB
		mock            sqlmock.Sqlmock
		err             error
		responseObj     AddEventApiResponse
		listResponseObj GetEventsApiResponse
		r               *gin.Engine
	)

	BeforeEach(func() {
		fakeDbClient = &fakeDb.FakeDBClient{}
		client = database.NewClient(fakeDbClient)
		service = services.NewClient(client)
		NewEventsApiClient(client, service)
		db, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())
		r = getRouter(false)
		r.POST("/add", AddEvent)
		r.GET("/events", GetEvents)
	})

	It("should add a record with no errors", func() {
		//Prepare all the mocks
		mock.ExpectQuery("SELECT id, name, description, date_added FROM event where date_added = ?").WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "date_added"}).AddRow(1, "test name", "test description", time.Now()))
		row := db.QueryRow("SELECT id, name, description, date_added FROM event where date_added = ?", 1)
		fakeDbClient.QueryRowCall.Returns.Row = row
		fakeDbClient.ExecCall.Returns.Result = fakeDb.NewPositiveResult(1000, 1, nil)

		byteArray, err := json.Marshal(gin.H{
			"name":        "Avenged Sevenfold Summer Tour",
			"description": "The sickest metal band is back!",
		})
		Expect(err).ToNot(HaveOccurred())

		//establish a request and fire it.
		request, err := http.NewRequest("POST", "/add", bytes.NewReader(byteArray))
		request.Header.Add("Content-Type", "application/json")

		Expect(err).ShouldNot(HaveOccurred())
		testHTTPResponse(r, request, func(writer *httptest.ResponseRecorder) bool {
			fmt.Println("code:", writer.Code)
			statusOK := writer.Code == http.StatusOK
			Expect(statusOK).To(BeTrue())
			byteArray, err := ioutil.ReadAll(writer.Body)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(byteArray, &responseObj)
			Expect(err).ToNot(HaveOccurred())
			return statusOK
		})
	})

	It("should reject an add with no POST body with a bad request", func() {
		reader := bytes.NewReader([]byte{})
		request, err := http.NewRequest("POST", "/add", reader)
		Expect(err).ToNot(HaveOccurred())
		testHTTPResponse(r, request, func(w *httptest.ResponseRecorder) bool {
			statusBadRequest := w.Code == http.StatusBadRequest
			Expect(statusBadRequest).To(BeTrue())
			return statusBadRequest
		})
	})

	It("should return the record it created", func() {
		Expect(responseObj.doesItPassTheContract()).To(BeTrue())
	})

	It("should get all records without errors", func() {
		rowsObj := sqlmock.NewRows([]string{"id", "name", "description"})
		for i := 0; i < 50; i++ {
			rowsObj.AddRow(i, fmt.Sprintf("test name %v", i), fmt.Sprintf("test desc %v", i))
		}
		mock.ExpectQuery("SELECT").WillReturnRows(rowsObj)
		rows, err := db.Query("SELECT")
		Expect(err).ToNot(HaveOccurred())
		fakeDbClient.QueryCall.Returns.Rows = rows
		fakeDbClient.QueryCall.Returns.Error = nil

		request, err := http.NewRequest("GET", "/events", nil)
		Expect(err).ToNot(HaveOccurred())
		var statusOK bool
		testHTTPResponse(r, request, func(w *httptest.ResponseRecorder) bool {
			statusOK = w.Code == http.StatusOK
			byteArray, err := ioutil.ReadAll(w.Body)
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal(byteArray, &listResponseObj)
			Expect(err).ToNot(HaveOccurred())
			return statusOK
		})
		Expect(statusOK).To(BeTrue())
		isntNil := assert.NotNil(nil, listResponseObj, nil)
		fmt.Println("responseObj isnt nil:", isntNil, responseObj)
	})
})
