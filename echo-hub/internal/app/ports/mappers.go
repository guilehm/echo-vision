package ports

import "github.com/guilehm/echo-vision/echo-hub/internal/app/domain"

func MapEventToApiResponse(e *domain.Event) *EventResponse {
	return &EventResponse{
		UserID:    e.UserID(),
		ID:        e.ID(),
		EventType: e.EventType().String(),
		SubType:   e.SubType().String(),
		Status:    e.Status().String(),
		CreatedAt: e.CreatedAt(),
		UpdateAt:  e.UpdatedAt(),
	}
}

func MapEventsToApiResponse(events []*domain.Event) []*EventResponse {
	apiResponseResults := make([]*EventResponse, 0, len(events))
	for i := range events {
		apiResponseResults = append(apiResponseResults, MapEventToApiResponse(events[i]))
	}
	return apiResponseResults
}

func MapUserToApiResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
		Email:     u.Email(),
	}
}
