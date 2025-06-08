package publishers

import (
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

type PublisherGroup struct {
	publisher messaging.Publisher
}

func NewPublisherGroup(publisher messaging.Publisher) ports.PublisherPort {
	return &PublisherGroup{
		publisher: publisher,
	}
}
