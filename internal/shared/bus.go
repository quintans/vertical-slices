package shared

import (
	"context"

	"github.com/quintans/vertical-slices/internal/lib/bus"
)

type Publisher interface {
	Publish(ctx context.Context, m ...bus.Message) error
}
