package tests

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
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

	Context("when processing an image analysis event", func() {
		It("should successfully persist the event and publish the event status update message", func() {
			// Arrange
			data := json.RawMessage(`{"labels": ["cat", "dog"]}`)

			statusUpdateCalled := make(chan any)
			topic := hubevents.EventImageAnalysisStatusUpdatedProcessing
			var msg hubevents.EventStatusUpdateMessage

			handler.Mock.On("Handle", mock.Anything, mock.Anything).Once().
				Return(messaging.Success).
				Run(func(args mock.Arguments) {
					defer GinkgoRecover()

					Expect(args).To(HaveLen(2))
					Expect(args.Get(0)).ToNot(BeNil())
					Expect(args.Get(1)).ToNot(BeNil())

					message := args.Get(1).(messaging.Message)
					Expect(message.Topic).To(Equal(topic))

					err := json.Unmarshal(message.Payload, &msg)
					Expect(err).ToNot(HaveOccurred())

					Expect(msg).ToNot(BeNil())
					Expect(msg.ID).To(Equal(eventID))
					Expect(msg.Status).To(Equal(hubevents.EventStatusProcessing))
					Expect(jsonEqual(msg.Data, data)).To(BeTrue())

					close(statusUpdateCalled)
				})

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
