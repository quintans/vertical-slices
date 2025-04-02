package bus

import (
	"context"
)

type Handler[T Message] func(context.Context, T) error

type Message interface {
	Kind() string
}

type Bus struct {
	handlers map[string][]Handler[Message]
}

func New() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler[Message]),
	}
}

func Register[T Message](bus *Bus, handler func(context.Context, T) error) {
	var zero T
	kind := zero.Kind()

	handlers := bus.handlers[kind]
	if handlers == nil {
		handlers = []Handler[Message]{}
	}

	handlers = append(handlers, func(ctx context.Context, m Message) error {
		return handler(ctx, m.(T))
	})

	bus.handlers[kind] = handlers
}

func (b *Bus) Publish(ctx context.Context, msgs ...Message) error {
	for k := range msgs {
		handlers := b.handlers[msgs[k].Kind()]

		for _, handler := range handlers {
			return handler(ctx, msgs[k])
		}
	}

	return nil
}
