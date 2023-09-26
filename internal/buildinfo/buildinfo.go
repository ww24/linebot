package buildinfo

var (
	// version is set during build
	version string
	// revision is set during build
	//nolint:gochecknoglobals
	revision string
	// repository is set during build
	//nolint:gochecknoglobals
	repository string
)

// Version returns version string.
func Version() string {
	return version
}

// Revision returns revision string.
func Revision() string {
	return revision
}

func Repository() string {
	return repository
}
