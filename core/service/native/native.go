package native

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dirty-bro-tech/peers-touch-go/core/client"
	"github.com/dirty-bro-tech/peers-touch-go/core/cmd"
	cfg "github.com/dirty-bro-tech/peers-touch-go/core/config"
	"github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/peers/config"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/core/service"
	"github.com/dirty-bro-tech/peers-touch-go/core/util/log"
)

type native struct {
	opts *service.Options

	once sync.Once
}

func (s *native) Name() string {
	return s.opts.Name
}

func (s *native) Options() *service.Options {
	return s.opts
}

func (s *native) Client() client.Client {
	return s.opts.Client
}

func (s *native) Server() server.Server {
	return s.opts.Server
}

func (s *native) String() string {
	return "peers"
}

func (s *native) Start(ctx context.Context) error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	ready := make(chan interface{})
	// Start server in a goroutine
	go func() {
		if err := s.opts.Server.Start(ctx, server.WithReadyChan(ready)); err != nil {
			logger.Errorf(ctx, "server start failed: %v", err)
			return
		}
	}()

	// Wait for server to be ready
	select {
	case info := <-ready:
		logger.Infof(ctx, "server started: %v", info)
	case <-ctx.Done():
		return ctx.Err()
	}

	// todo start achievement, but now it's blocked by Start
	// so will never be called. it should be pass a exit signal
	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	logger.Infof(ctx, "peers' service started at %s", s.opts.Server.Options().Address)
	return nil
}

func (s *native) Stop(ctx context.Context) error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(ctx); err != nil {
		return err
	}

	if err := s.opts.Config.Close(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *native) Run() error {
	if err := s.Start(s.opts.Ctx()); err != nil {
		return err
	}

	logger.Infof(s.opts.Ctx(), "[%s] server started", s.Name())

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	}

	select {
	// wait on kill signal
	case <-ch:
		logger.Warnf(s.opts.Ctx(), "received signal, stopping")
	case <-s.opts.Ctx().Done():
		// todo, try to store canceled context data for logger
		logger.Infof(context.Background(), "[%s] server stopped", s.Name())
	}

	return s.Stop(s.opts.Ctx())
}

// NewService
// todo remove rootOptions, not graceful
func NewService(rootOpts *option.Options, opts ...option.Option) service.Service {
	defaultOpts := []option.Option{
		// todo remove non-peers' code
		service.RPC("peers"),
		// load config
		service.BeforeInit(func(sOpts *service.Options) error {
			// cmd helps peers parse command options and reset the options that should work.
			err := sOpts.Cmd.Init()
			if err != nil {
				log.Errorf("cmd init error: %s", err)
				return err
			}

			err = config.LoadConfig(sOpts)
			if err != nil {
				log.Errorf("load config error: %s", err)
				return err
			}

			return nil
		}),
		// parse config to options for components
		service.BeforeInit(func(sOpts *service.Options) error {
			err := config.SetOptions(sOpts)
			if err != nil {
				log.Errorf("init components' options error: %s", err)
				return err
			}

			return nil
		}),

		// set the default components
		// see initComponents
		// service.Store(plugin.StorePlugins["native"].New()),
		// service.Server(plugin.ServerPlugins["native"].New()),
		// service.Logger(plugin.LoggerPlugins["console"].New()),
		service.Config(cfg.DefaultConfig),
		//service.HandleSignal(true),
	}

	// make sure the custom options are applied last
	// so that they can override the default options
	defaultOpts = append(defaultOpts, opts...)

	for _, o := range defaultOpts {
		rootOpts.Apply(o)
	}

	s := &native{
		opts: service.GetOptions(rootOpts),
	}

	// set CMD for loading config
	s.opts.Cmd = cmd.NewCmd()

	return s
}
