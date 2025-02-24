package server

// region server options

// Option is a function that can be used to configure a server
type Option func(*Options)

// Options is the server options
type Options struct {
	Address  string            `pconf:"address"` // Server address
	Timeout  int               `pconf:"timout"`  // Server timeout
	Metadata map[string]string `pconf:"metadata"`
	Handlers []Handler
}

// WithAddress sets the server address
func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithTimeout sets the server timeout
func WithTimeout(timeout int) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithMetadata associated with the server
func WithMetadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

func WithHandlers(handlers ...Handler) Option {
	return func(o *Options) {
		o.Handlers = handlers
	}
}

// endregion
