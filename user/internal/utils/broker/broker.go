package broker

type Producer interface {
	Produce(value ...any) error
}

type Consumer interface {
	Consume() any
}
