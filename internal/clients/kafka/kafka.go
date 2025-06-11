package kafka

type Consumer interface {
	Start() error
	Stop() error
}
