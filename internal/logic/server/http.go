package server

import (
	"context"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/api/fileupload"
	"github.com/ykds/zura/internal/logic/api/friend_application"
	"github.com/ykds/zura/internal/logic/api/friends"
	"github.com/ykds/zura/internal/logic/api/group"
	"github.com/ykds/zura/internal/logic/api/message"
	"github.com/ykds/zura/internal/logic/api/session"
	"github.com/ykds/zura/internal/logic/api/user"
	"github.com/ykds/zura/internal/logic/config"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/pprof"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*gin.Engine
	c          config.HttpServerConfig
	l          log.Logger
	httpServer *http.Server
	debug      bool
}

type Option func(*HttpServer)

func WithLogger(l log.Logger) Option {
	return func(hs *HttpServer) {
		hs.l = l
	}
}

func WithDebug(debug bool) Option {
	return func(hs *HttpServer) {
		hs.debug = debug
	}
}

func NewHttpServer(cfg config.HttpServerConfig, opts ...Option) *HttpServer {
	server := &HttpServer{
		c: cfg,
	}
	for _, opt := range opts {
		opt(server)
	}

	var engine *gin.Engine
	if server.debug {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.LoggerWithWriter(log.GetGlobalLogger()), gin.RecoveryWithWriter(log.GetGlobalLogger()))
	}
	engine.Static(common.StaticPath, common.StaticDir)
	pprof.RouteRegister(engine)
	loadRouters(engine)
	server.httpServer = &http.Server{
		Addr:    ":" + server.c.Port,
		Handler: engine,
	}
	return server
}

func (h *HttpServer) Run() {
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			if h.l != nil {
				h.l.Fatalf("http server exit, error: %+v", err)
			}
		}
	}()
}

func (h *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	h.httpServer.Shutdown(ctx)
}

func loadRouters(r gin.IRouter) {
	api := r.Group("/api")
	user.RegisterUserRouter(api)
	friends.RegisterFriendsRouter(api)
	friend_application.RegisterFriendApplicationRouter(api)
	session.RegisterSessionRouter(api)
	fileupload.RegisterUploadRouter(api)
	group.RegisterGroupRouter(api)
	message.RegisterMessageRouter(api)
}
