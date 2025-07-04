package ports

import (
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

func MapEventToApiResponse(e *domain.Event) *dtos.EventResponse {
	var fileResponse *dtos.FileResponse
	if e.File() != nil {
		fileResponse = &dtos.FileResponse{
			Filename:    e.File().Filename,
			Filepath:    e.File().Filepath,
			ContentType: e.File().ContentType,
			Filesize:    e.File().Filesize,
		}
	}

	return &dtos.EventResponse{
		UserID:    e.UserID(),
		ID:        e.ID(),
		EventType: e.EventType().String(),
		SubType:   e.SubType().String(),
		Status:    e.Status().String(),
		Result:    e.Result(),
		File:      fileResponse,
		CreatedAt: e.CreatedAt(),
		UpdateAt:  e.UpdatedAt(),
	}
}

func MapEventsToApiResponse(events []*domain.Event) []*dtos.EventResponse {
	apiResponseResults := make([]*dtos.EventResponse, 0, len(events))
	for i := range events {
		apiResponseResults = append(apiResponseResults, MapEventToApiResponse(events[i]))
	}
	return apiResponseResults
}

func MapUserToApiResponse(u *domain.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:        u.ID(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
		Email:     u.Email(),
	}
}

func MapEventToMessage(e *domain.Event) ([]byte, error) {
	message := &hubevents.EventMessage{
		ID:      e.ID(),
		Type:    e.EventType(),
		SubType: e.SubType(),
		File:    e.File(),
	}
	return message.ToJSON()
}
