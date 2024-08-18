package log

type options struct {
	service      Service
	version      Version
	repository   Repository
	revisionID   RevisionID
	gcpProjectID GCPProjectID
}

func (o *options) apply(opts []Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

type Option interface {
	apply(*options)
}

type Service string

func (s Service) apply(opts *options) {
	opts.service = s
}

type Version string

func (v Version) apply(opts *options) {
	opts.version = v
}

type Repository string

func (r Repository) apply(opts *options) {
	opts.repository = r
}

type RevisionID string

func (r RevisionID) apply(opts *options) {
	opts.revisionID = r
}

type GCPProjectID string

func (p GCPProjectID) apply(opts *options) {
	opts.gcpProjectID = p
}
