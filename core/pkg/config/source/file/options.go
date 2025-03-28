package file

import (
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/source"
)

type filePathKey struct{}

// WithPath sets the path to file
func WithPath(p string) option.Option {
	return source.WrapOption(func(o *source.Options) {
		o.AppendCtx(filePathKey{}, p)
	})
}
