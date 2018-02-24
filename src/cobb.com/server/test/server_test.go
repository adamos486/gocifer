package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//. "cobb.com/server"

	"cobb.com/server/database"
	eventController "cobb.com/server/events/controllers"
	fakeDb "cobb.com/server/fakes/database"
	"cobb.com/server/events/services"
)

var _ = Describe("Server", func() {
	RegisterFailHandler(Fail)

	var (
		client       database.Client
		fakeDbClient *fakeDb.FakeDBClient
		service      *services.EventsServiceClient
	)

	BeforeEach(func() {
		fakeDbClient = &fakeDb.FakeDBClient{}
		client = database.NewClient(fakeDbClient)
		service = services.NewClient(client)
		eventController.NewEventsApiClient(service)
	})

	It("should be able to add an Event", func() {
		return
	})
	It("should have a properly formed response", func() {
		return
	})
})
