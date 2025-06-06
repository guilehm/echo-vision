package analysistypes

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
