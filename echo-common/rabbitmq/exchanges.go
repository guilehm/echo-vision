package rabbitmq

type ExchangeName string

func (e ExchangeName) String() string {
	return string(e)
}
