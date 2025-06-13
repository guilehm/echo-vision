package messaging

type ExchangeType string

func (et ExchangeType) String() string {
	return string(et)
}

const (
	ExchangeTypeDirect  ExchangeType = "direct"
	ExchangeTypeFanout  ExchangeType = "fanout"
	ExchangeTypeTopic   ExchangeType = "topic"
	ExchangeTypeHeaders ExchangeType = "headers"
)

// Message defines the structure of a message handled by the async messaging system.
type Message struct {
	Topic   string
	Payload []byte
	Headers map[string]string
}
