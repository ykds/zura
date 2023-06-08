package server

import (
	"context"
	"net/http"
	"time"
	"zura/internal/api/user"
	"zura/pkg/log"

	"github.com/gin-gonic/gin"
)

type HttpServerConfig struct {
	Port string `json:"port"`
}

func DefaultConfig() HttpServerConfig {
	return HttpServerConfig{
		Port: "8000",
	}
}

type HttpServer struct {
	*gin.Engine
	ctx        context.Context
	cancel     context.CancelFunc
	c          HttpServerConfig
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

func WithConfig(c HttpServerConfig) Option {
	return func(hs *HttpServer) {
		hs.c = c
	}
}

func NewHttpServer(ctx context.Context, opts ...Option) *HttpServer {
	ctx2, cancel := context.WithCancel(ctx)
	server := &HttpServer{
		ctx:    ctx2,
		cancel: cancel,
		c:      DefaultConfig(),
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
	go func() {
		select {
		case <-h.ctx.Done():
			h.Shutdown()
			if h.l != nil {
				h.l.Info("stop http server")
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
}
