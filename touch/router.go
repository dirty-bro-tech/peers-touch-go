package touch

import (
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/touch/model"
)

const (
	RoutersNameManagement  = "management"
	RoutersNameActivityPub = "activitypub"
	RoutersNameWellKnown   = ".well-known"
	RoutersNameUser        = "user"
	RoutersNamePeer        = "peer"
)

// Router is a server handler that can be registered with a server.
// Peers defines a router protocol that can be used to register handlers with a server.
// also supplies standard handlers which follow activityPub protocol.
// if you what to register a handler with Peers server, you can implement this interface, then call server.listPeers() to register it.
type Router server.Handler

type Routers interface {
	Routers() []Router

	// Name declares the cluster-name of routers
	// it must be unique. Peers uses it to check if there are already routers(like activitypub
	// and management interface.) that must be registered,
	// if you want to register a bundle of routers with the same name, it will be overwritten
	Name() string
}

type RouterURL string

func (apr RouterURL) Name() string {
	return string(apr)
}

func (apr RouterURL) URL() string {
	return string(apr)
}

func Handlers() []option.Option {
	m := NewManageRouter()
	a := NewActivityPubRouter()
	w := NewWellKnownRouter()
	u := NewUserRouter()
	p := NewPeerRouter()

	handlers := make([]option.Option, 0)

	for _, r := range m.Routers() {
		handlers = append(handlers, server.WithHandlers(convertRouterToServerHandler(r)))
	}

	for _, r := range a.Routers() {
		handlers = append(handlers, server.WithHandlers(convertRouterToServerHandler(r)))
	}

	for _, r := range w.Routers() {
		handlers = append(handlers, server.WithHandlers(convertRouterToServerHandler(r)))
	}

	for _, r := range u.Routers() {
		handlers = append(handlers, server.WithHandlers(convertRouterToServerHandler(r)))
	}

	for _, r := range p.Routers() {
		handlers = append(handlers, server.WithHandlers(convertRouterToServerHandler(r)))
	}

	return handlers
}

func convertRouterToServerHandler(r Router) server.Handler {
	return server.Handler(r)
}

func SuccessResponse(ctx *app.RequestContext, msg string, data interface{}) {
	if msg == "" {
		msg = "success"
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(msg, data))
}

func FailedResponse(ctx *app.RequestContext, err error) {
	if err != nil {
		var e *model.Error
		if errors.As(err, &e) {
			ctx.JSON(http.StatusBadRequest, e)
			return
		}

		ctx.JSON(http.StatusBadRequest, model.UndefinedError(err))
		return
	}

	ctx.JSON(http.StatusBadRequest, model.ErrUndefined)
}
