package pprof

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

func RouteRegister(rg gin.IRouter) {
	r := rg.Group("/debug/pprof")
	{
		r.GET("/", gin.WrapF(pprof.Index))
		r.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		r.GET("/profile", gin.WrapF(pprof.Profile))
		r.GET("/trace", gin.WrapF(pprof.Trace))
		r.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		r.GET("/block", gin.WrapH(pprof.Handler("block")))
		r.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		r.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		r.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		r.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
