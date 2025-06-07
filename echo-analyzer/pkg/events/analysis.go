package analyzerevents

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// EventStatus represents the status of an event.
type EventStatus string

// EventStatus values represent the possible states of an event.
const (
	EventStatusProcessing EventStatus = "processing"
	EventStatusCompleted  EventStatus = "completed"
	EventStatusFailed     EventStatus = "failed"
)

// Values returns all possible EventStatus values.
func (es EventStatus) Values() []EventStatus {
	return []EventStatus{
		EventStatusProcessing,
		EventStatusCompleted,
		EventStatusFailed,
	}
}

// String returns the string representation of the EventStatus.
func (es EventStatus) String() string {
	return string(es)
}

// EventStatusUpdateMessage represents a message for updating the status of an image analysis event.
const (
	EventImageAnalysisStatusUpdatedGeneric    = "analyzer.event.image_analysis.status_updated.*"
	EventImageAnalysisStatusUpdatedProcessing = "analyzer.event.image_analysis.status_updated.processing"
	EventImageAnalysisStatusUpdatedCompleted  = "analyzer.event.image_analysis.status_updated.completed"
	EventImageAnalysisStatusUpdatedFailed     = "analyzer.event.image_analysis.status_updated.failed"
)

// EventImageAnalysisStatusUpdatedTopic builds a topic string for image analysis status updates.
func BuildEventImageAnalysisStatusUpdatedTopic(status EventStatus) string {
	return fmt.Sprintf("analyzer.event.image_analysis.status_updated.%s", status)
}

// EventImageAnalysisStatusUpdateMessage represents a message for updating the status of an image analysis event.
type EventImageAnalysisStatusUpdateMessage struct {
	ID     uuid.UUID       `json:"id"`
	Status EventStatus     `json:"status"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// VisionAnalysisResult represents the result of an image analysis.
type VisionAnalysisResult struct {
	Labels []Label `json:"labels"`
}

// Label represents an identified label or concept in the image.
type Label struct {
	Aliases    []string   `json:"aliases"`
	Categories []string   `json:"categories"`
	Confidence *float32   `json:"confidence"`
	Instances  []Instance `json:"instances"`
	Name       *string    `json:"name"`
	Parents    []string   `json:"parents"`
}

// Instance represents an instance of a label detected in the image.
type Instance struct {
	BoundingBox BoundingBox `json:"boundingBox"`
	Confidence  *float32    `json:"confidence"`
}

// FaceDetail represents details about a detected face.
type FaceDetail struct {
	BoundingBox BoundingBox `json:"boundingBox"`
	Confidence  *float32    `json:"confidence"`
	Emotions    []Emotion   `json:"emotions"`
}

// BoundingBox represents the dimensions and position of an object in the image.
type BoundingBox struct {
	Height *float32 `json:"height"`
	Width  *float32 `json:"width"`
	Top    *float32 `json:"top"`
	Left   *float32 `json:"left"`
}

type Emotion struct {
	Type       string
	Confidence *float32
}
