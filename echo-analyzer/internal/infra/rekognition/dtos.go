package awsrekognition

import (
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	analysistypes "github.com/guilehm/echo-vision/echo-analyzer/pkg/types"
)

func labelToDomain(
	label types.Label,
) analysistypes.Label {
	as := make([]string, 0, len(label.Aliases))
	for _, a := range label.Aliases {
		if a.Name != nil {
			as = append(as, *a.Name)
		}
	}

	cs := make([]string, 0, len(label.Categories))
	for _, c := range label.Categories {
		if c.Name != nil {
			cs = append(cs, *c.Name)
		}
	}

	is := make([]analysistypes.Instance, 0, len(label.Instances))
	for _, i := range label.Instances {
		is = append(is, analysistypes.Instance{
			Confidence: i.Confidence,
			BoundingBox: analysistypes.BoundingBox{
				Height: i.BoundingBox.Height,
				Left:   i.BoundingBox.Left,
				Top:    i.BoundingBox.Top,
				Width:  i.BoundingBox.Width,
			},
		})
	}

	ps := make([]string, 0, len(label.Parents))
	for _, p := range label.Parents {
		if p.Name != nil {
			ps = append(ps, *p.Name)
		}
	}

	return analysistypes.Label{
		Aliases:    as,
		Categories: cs,
		Confidence: label.Confidence,
		Instances:  is,
		Name:       label.Name,
		Parents:    ps,
	}
}
