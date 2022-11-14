package domain

import "fmt"

type OnMessageReceive func(string)

type EventSubscriberService interface {
	OnMessageReceive(value string)
}

type EventSubscriberServiceImpl struct{}

func NewEventSubscriberServiceImpl() EventSubscriberService {
	return &EventSubscriberServiceImpl{}
}

func (subscriberService *EventSubscriberServiceImpl) OnMessageReceive(value string) {
	fmt.Println("received at callback: ", value)
}
