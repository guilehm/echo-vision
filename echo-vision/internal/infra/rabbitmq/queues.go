package rabbitmq

type QueueName string

func (q QueueName) String() string {
	return string(q)
}
