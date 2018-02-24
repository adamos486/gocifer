package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"cobb.com/server/database"
	. "cobb.com/server/events/services"
	fakeDB "cobb.com/server/fakes/database"
	"cobb.com/server/utils"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"time"
	//"fmt"
	"cobb.com/server/events/models"
	"fmt"
	"reflect"
)

var _ = Describe("Event Services Tests", func() {
	//Create fakes here.
	var (
		client     database.Client
		fakeCollab *fakeDB.FakeDBClient
		subject    *EventsServiceClient
		db         *sql.DB
		mock       sqlmock.Sqlmock
		err        error
	)

	BeforeEach(func() {
		//Initialize Client definition
		fakeCollab = &fakeDB.FakeDBClient{}
		client = database.NewClient(fakeCollab)
		subject = NewClient(client)
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("AddCannedRowToEvents DB", func() {
		It("should insert a canned row into the event database", func() {
			//Canned fakes go here.
			fakeCollab.ExecCall.Returns.Result = fakeDB.NewPositiveResult(1, 1, nil)

			//Test method in question
			lookup, result, err := subject.AddCannedRowToEventsDB()

			//Validate expected result
			Expect(err).ToNot(HaveOccurred())
			assert.NotNil(nil, lookup, "")
			assert.NotNil(nil, result, "")
		})

		It("should pass any error back to the caller", func() {
			//Canned fakes go here.
			fakeCollab.ExecCall.Returns.Error = errors.New("totally real not test error")

			//Test method in question
			lookup, result, err := subject.AddCannedRowToEventsDB()

			//Validate expected result
			Expect(err).To(HaveOccurred())
			assert.Nil(nil, lookup, "")
			assert.Nil(nil, result, "")
		})
	})

	var _ = Describe("AddNewEvent DB", func() {
		It("should insert a row into the event database with passed in info", func() {
			fakeCollab.ExecCall.Returns.Result = fakeDB.NewPositiveResult(2, 1, nil)

			now := time.Now()
			mockRows := sqlmock.NewRows([]string{"id", "name", "description", "date_added"})
			mockRows = mockRows.AddRow(0, "test concert", "test description", now)
			mock.ExpectQuery("SELECT").WillReturnRows(mockRows)
			row := db.QueryRow("SELECT")
			fakeCollab.QueryRowCall.Returns.Row = row

			comparison := models.EventRow{
				ID:          0,
				Name:        "test concert",
				Description: "test description",
				DateAdded:   now,
			}

			//Test method
			created, result, err := subject.AddNewEvent("test concert", "test description")

			//Validate expected result
			Expect(err).ToNot(HaveOccurred())

			assert.NotNil(nil, created, "")
			assert.NotNil(nil, result, "")

			isEmpty := utils.IsEmpty(created)
			Expect(isEmpty).ToNot(BeTrue())

			//Make sure it's the same row.
			Expect(*created).To(Equal(comparison))

			isEmpty = utils.IsEmpty(result)
			Expect(isEmpty).ToNot(BeTrue())
		})

		It("should pass any error back to the caller", func() {
			fakeCollab.ExecCall.Returns.Error = errors.New("any error really")

			//Test method
			created, result, err := subject.AddNewEvent("test concert", "test description")

			//Test assumptions
			Expect(err).To(HaveOccurred())

			assert.Nil(nil, created, "")
			assert.Nil(nil, result, "")
		})

		It("should pass back specific error if row is nil and result fails", func() {
			fakeCollab.ExecCall.Returns.Result = fakeDB.NewPositiveResult(0, 0, nil)

			created, result, err := subject.AddNewEvent("test concert", "test description")

			assert.Nil(nil, created, "")
			assert.Nil(nil, result, "")
			assert.NotNil(nil, err, "")

			Expect(err.Error()).To(Equal("404: Event not found"))
		})

		It("should pass back a different error if row is nil and result passes", func() {
			fakeCollab.ExecCall.Returns.Result = fakeDB.NewPositiveResult(0, 1, nil)

			created, result, err := subject.AddNewEvent("test concert", "test description")

			assert.Nil(nil, created, "")
			assert.Nil(nil, result, "")
			assert.NotNil(nil, err, "")

			Expect(err.Error()).To(Equal("500: Something is broken in AddNewEvent fetching"))
		})
	})

	var _ = Describe("GetEvents From DB", func() {
		It("should get records without an error", func() {
			fakeRows := sqlmock.NewRows([]string{"id", "name", "description"})
			for i := 0; i < 50; i++ {
				fakeRows.AddRow(i, fmt.Sprintf("test name %v", i), fmt.Sprintf("test desc %v", i))
			}
			mock.ExpectQuery("SELECT id, name, description FROM event;").WillReturnRows(fakeRows)
			rows, err := db.Query("SELECT id, name, description FROM event;")
			Expect(err).ToNot(HaveOccurred())

			fakeCollab.QueryCall.Returns.Rows = rows

			resultRows, err := subject.GetAllEvents()
			Expect(err).ToNot(HaveOccurred())
			Expect(*resultRows).ToNot(BeEmpty())
			Expect(len(*resultRows)).To(Equal(50))
		})

		It("should get back an array of GetAllEventsResponseStructs", func() {
			fakeRows := sqlmock.NewRows([]string{"id", "name", "description"})
			fakeRows.AddRow(1, "test name", "test desc")
			mock.ExpectQuery("SELECT id, name, description FROM event;").WillReturnRows(fakeRows)
			rows, err := db.Query("SELECT id, name, description FROM event;")
			Expect(err).ToNot(HaveOccurred())

			fakeCollab.QueryCall.Returns.Rows = rows

			resultRows, err := subject.GetAllEvents()
			Expect(err).ToNot(HaveOccurred())
			Expect(*resultRows).ToNot(BeEmpty())

			singleRow := (*resultRows)[0]
			t := reflect.TypeOf(singleRow)
			Expect(t.Name()).To(Equal("GetAllEventsResponseStruct"))
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				Expect(f.Name).ToNot(Equal("DateAdded"))
			}
		})

		It("should return an empty array if no records", func() {
			fakeRows := sqlmock.NewRows([]string{"id", "name", "description"})
			mock.ExpectQuery("SELECT id, name, description FROM event;").WillReturnRows(fakeRows)
			rows, err := db.Query("SELECT id, name, description FROM event;")
			Expect(err).ToNot(HaveOccurred())

			fakeCollab.QueryCall.Returns.Rows = rows

			resultRows, err := subject.GetAllEvents()

			Expect(err).ToNot(HaveOccurred())
			Expect(*resultRows).ToNot(BeNil())
			Expect(*resultRows).To(BeEmpty())
		})

		It("should pass all errors back to caller", func() {
			fakeCollab.QueryCall.Returns.Error = errors.New("i'm a little teapot")

			resultRows, err := subject.GetAllEvents()

			Expect(err).ToNot(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(resultRows).To(BeNil())
		})
	})
})
