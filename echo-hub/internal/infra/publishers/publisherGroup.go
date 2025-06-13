package publishers

import (
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

type PublisherGroup struct {
	publisher messaging.Publisher
}

func NewPublisherGroup(publisher messaging.Publisher) ports.PublisherPort {
	return &PublisherGroup{
		publisher: publisher,
	}
}
