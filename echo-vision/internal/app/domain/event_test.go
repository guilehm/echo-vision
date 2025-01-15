package domain_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

func TestEvent_Validate(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		eventType   domain.EventType
		subType     domain.EventSubType
		status      domain.EventStatus
		payload     json.RawMessage
		result      json.RawMessage
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "Valid event",
			id:          uuid.New(),
			eventType:   domain.EventTypeImageAnalysis,
			subType:     domain.EventSubTypeDetectLabels,
			status:      domain.EventStatusPending,
			payload:     json.RawMessage(`{"key": "value"}`),
			result:      json.RawMessage(`{"result": "success"}`),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Invalid event ID",
			id:          uuid.Nil,
			eventType:   domain.EventTypeImageAnalysis,
			subType:     domain.EventSubTypeDetectLabels,
			status:      domain.EventStatusPending,
			payload:     json.RawMessage(`{"key": "value"}`),
			result:      json.RawMessage(`{"result": "success"}`),
			wantErr:     true,
			expectedErr: shared.ErrInvalidID,
		},
		{
			name:        "Invalid event type",
			id:          uuid.New(),
			eventType:   "",
			subType:     domain.EventSubTypeDetectLabels,
			status:      domain.EventStatusPending,
			payload:     json.RawMessage(`{"key": "value"}`),
			result:      json.RawMessage(`{"result": "success"}`),
			wantErr:     true,
			expectedErr: shared.ErrInvalidEventType,
		},
		{
			name:        "Nil payload",
			id:          uuid.New(),
			eventType:   domain.EventTypeImageAnalysis,
			subType:     domain.EventSubTypeDetectLabels,
			status:      domain.EventStatusPending,
			payload:     nil,
			result:      json.RawMessage(`{"result": "updated"}`),
			wantErr:     true,
			expectedErr: shared.ErrInvalidPayload,
		},
		{
			name:        "Empty status",
			id:          uuid.New(),
			eventType:   domain.EventTypeImageAnalysis,
			subType:     domain.EventSubTypeDetectLabels,
			status:      "",
			payload:     json.RawMessage(`{"key": "value"}`),
			result:      json.RawMessage(`{"result": "deleted"}`),
			wantErr:     true,
			expectedErr: shared.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := domain.NewEvent(tt.id, tt.eventType, tt.subType, tt.payload, tt.result, tt.status)

			gotErr := e.Validate()
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", gotErr, tt.wantErr)
			}

			if gotErr != nil && !errors.Is(gotErr, tt.expectedErr) {
				t.Errorf("Validate() error type = %v, expected %v", gotErr, tt.expectedErr)
			}
		})
	}
}
