package reader

import (
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder/hcl"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder/json"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder/toml"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder/xml"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/encoder/yaml"
)

type Options struct {
	Encoding map[string]encoder.Encoder
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]encoder.Encoder{
			"json": json.NewEncoder(),
			"yaml": yaml.NewEncoder(),
			"toml": toml.NewEncoder(),
			"xml":  xml.NewEncoder(),
			"hcl":  hcl.NewEncoder(),
			"yml":  yaml.NewEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		if o.Encoding == nil {
			o.Encoding = make(map[string]encoder.Encoder)
		}
		o.Encoding[e.String()] = e
	}
}
