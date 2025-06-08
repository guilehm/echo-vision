package tests

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image Analysis Consumer", func() {
	var u *domain.User
	eventID := uuid.MustParse("82f86ed8-df21-4587-9de0-d8b8bffeafd7")

	BeforeEach(func() {
		u = saveUser(makeUser("arthur@gmail.com"))
		event := domain.NewEvent(
			u.ID(),
			eventID,
			hubevents.EventTypeImageAnalysis,
			hubevents.EventSubTypeDetectLabels,
			nil,
			hubevents.EventStatusPending,
			nil,
			time.Now(),
			time.Now(),
		)
		_, err := repo.SaveEvent(ctx, nil, event)
		Expect(err).ToNot(HaveOccurred())
	})

	FContext("when processing an image analysis event", func() {
		It("should successfully persist the event and publish the event status update message", func() {
			// Arrange
			data := json.RawMessage(`{"labels": ["cat", "dog"]}`)

			statusUpdateCalled := make(chan any)
			topic := hubevents.EventImageAnalysisStatusUpdatedProcessing
			// var msg hubevents.EventStatusUpdateMessage

			expectFunc := func(msg hubevents.EventStatusUpdateMessage) {
				Expect(msg).ToNot(BeNil())
				Expect(msg.ID).To(Equal(eventID))
				Expect(msg.Status).To(Equal(hubevents.EventStatusProcessing))
				Expect(jsonEqual(msg.Data, data)).To(BeTrue())
			}

			// Expect the status update message to be handled
			expectMessageCalled(
				topic,
				statusUpdateCalled,
				expectFunc,
			)

			// Act
			handleMessage(
				analyzerevents.EventImageAnalysisStatusUpdatedProcessing,
				analyzerevents.EventImageAnalysisStatusUpdateMessage{
					ID:     eventID,
					Status: analyzerevents.EventStatusProcessing,
					Data:   data,
				},
			)

			// Assert
			Eventually(statusUpdateCalled, "1s").Should(BeClosed())
			updatedEvent, err := eventUseCase.FindEventByID(ctx, eventID)
			Expect(err).ToNot(HaveOccurred())

			Expect(updatedEvent).ToNot(BeNil())
			Expect(updatedEvent.ID()).To(Equal(eventID))
			Expect(updatedEvent.Status()).To(Equal(hubevents.EventStatusProcessing))
			Expect(updatedEvent.Result()).To(Equal(data))
		})
	})
})
