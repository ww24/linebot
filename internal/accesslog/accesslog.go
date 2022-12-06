package accesslog

import "github.com/google/wire"

// Set provides a wire set.
var Set = wire.NewSet(
	NewPublisher,
)
