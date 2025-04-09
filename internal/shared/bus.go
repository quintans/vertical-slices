package shared

import (
	"context"

	"github.com/quintans/vertical-slices/internal/lib/bus"
)

// Publisher publishes a message
// This declared in the shared package because it will be used accros all slices
type Publisher interface {
	Publish(ctx context.Context, m ...bus.Message) error
}
