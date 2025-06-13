package rabbitmqmocks

type MockedEvent struct {
	Topic   string
	Payload []byte
	Headers map[string]string
}
