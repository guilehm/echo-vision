package domain

// VisionAnalysisResult represents the result of an image analysis.
type VisionAnalysisResult struct {
	Labels []Label `json:"labels"`
}

// Label represents an identified label or concept in the image.
type Label struct {
	Aliases    []string   `json:"aliases"`
	Categories []string   `json:"categories"`
	Confidence float64    `json:"confidence"`
	Instances  []Instance `json:"instances"`
	Name       string     `json:"name"`
	Parents    []string   `json:"parents"`
}

// Instance represents an instance of a label detected in the image.
type Instance struct {
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
}

// FaceDetail represents details about a detected face.
type FaceDetail struct {
	BoundingBox BoundingBox `json:"bounding_box"`
	Confidence  float64     `json:"confidence"`
}

// BoundingBox represents the dimensions and position of an object in the image.
type BoundingBox struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Top    float64 `json:"top"`
	Left   float64 `json:"left"`
}
