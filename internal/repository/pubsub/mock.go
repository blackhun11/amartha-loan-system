package pubsub

import (
	"context"
	"fmt"
)

type Mock interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type mock struct{}

func NewMock() Mock {
	return &mock{}
}

func (m *mock) Publish(ctx context.Context, topic string, data []byte) error {
	// mock publish
	fmt.Println("mock publish to topic:", topic, "data:", string(data))
	return nil
}
