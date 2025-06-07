package tests

import (
	"encoding/json"
	"time"

	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Domain Validation", func() {
	var (
		validID   = uuid.New()
		validTime = time.Now()
	)

	eventCreationTests := []struct {
		name        string
		id          uuid.UUID
		eventType   hubevents.EventType
		subType     hubevents.EventSubType
		status      hubevents.EventStatus
		result      json.RawMessage
		wantErr     bool
		expectedErr error
		createdAt   time.Time
		updatedAt   time.Time
	}{
		{
			name:        "Valid event",
			id:          validID,
			eventType:   hubevents.EventTypeImageAnalysis,
			subType:     hubevents.EventSubTypeDetectLabels,
			status:      hubevents.EventStatusPending,
			result:      json.RawMessage(`{"result": "success"}`),
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Invalid event ID",
			id:          uuid.Nil,
			eventType:   hubevents.EventTypeImageAnalysis,
			subType:     hubevents.EventSubTypeDetectLabels,
			status:      hubevents.EventStatusPending,
			result:      json.RawMessage(`{"result": "success"}`),
			wantErr:     true,
			expectedErr: shared.ErrInvalidID,
		},
		{
			name:        "Invalid event type",
			id:          validID,
			eventType:   "",
			subType:     hubevents.EventSubTypeDetectLabels,
			status:      hubevents.EventStatusPending,
			result:      json.RawMessage(`{"result": "success"}`),
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     true,
			expectedErr: shared.ErrInvalidEventType,
		},
		// TODO: enable this test after implementing the payload validation
		// {
		// 	name:        "Nil payload",
		// 	id:          validID,
		// 	eventType:   hubevents.EventTypeImageAnalysis,
		// 	subType:     hubevents.EventSubTypeDetectLabels,
		// 	status:      hubevents.EventStatusPending,
		// 	payload:     nil,
		// 	result:      json.RawMessage(`{"result": "updated"}`),
		// 	createdAt:   validTime,
		// 	updatedAt:   validTime,
		// 	wantErr:     true,
		// 	expectedErr: shared.ErrInvalidPayload,
		// },
		{
			name:        "Empty status",
			id:          validID,
			eventType:   hubevents.EventTypeImageAnalysis,
			subType:     hubevents.EventSubTypeDetectLabels,
			status:      "",
			result:      json.RawMessage(`{"result": "deleted"}`),
			createdAt:   validTime,
			updatedAt:   validTime,
			wantErr:     true,
			expectedErr: shared.ErrInvalidStatus,
		},
	}

	Context("EventCreation", func() {
		for _, t := range eventCreationTests {
			It(t.name, func() {
				e := domain.NewEvent(
					validUser.ID(),
					t.id,
					t.eventType,
					t.subType,
					t.result,
					t.status,
					nil,
					t.createdAt,
					t.updatedAt,
				)

				gotErr := e.Validate()

				if t.wantErr {
					Expect(gotErr).To(HaveOccurred())
					Expect(gotErr).To(MatchError(t.expectedErr))
				} else {
					Expect(gotErr).NotTo(HaveOccurred())
				}
			})
		}
	})
})
