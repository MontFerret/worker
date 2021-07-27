package caching

type (
	Options struct {
		Size uint
	}

	Option func(opts *Options)
)

func NewOptions(setters ...Option) *Options {
	opts := new(Options)

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithSize(size uint) Option {
	return func(opts *Options) {
		opts.Size = size
	}
}
