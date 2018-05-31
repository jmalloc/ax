package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Validator wraps methods required to
// perform message validation
type Validator interface {
	Validate(ctx context.Context, msg ax.Message) error
}
