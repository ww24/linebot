package config

import "github.com/google/wire"

// Set provides a wire set.
var Set = wire.NewSet(
	NewLINEBot,
	NewScreenshot,
	NewStorage,
	NewTime,
	NewOtel,
	NewAccessLog,
	NewServiceEndpoint,
)
