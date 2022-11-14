package domain

type EventNotifierService interface {
	Send(value []byte) error
}
