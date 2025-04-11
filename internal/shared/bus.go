package shared

import (
	"context"

	"github.com/quintans/vertical-slices/internal/lib/eventbus"
)

// Publisher publishes a message
// This is declared in the shared package because it will be used accros all slices
type Publisher interface {
	Publish(ctx context.Context, m ...eventbus.Message) error
}
