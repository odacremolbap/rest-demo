package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"

	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/server/services"
)

// Server HTTP handling server
type Server struct {
	Port            int
	ShutDownTimeout time.Duration
}

// NewServer creates a new HTTP server
func NewServer(port int, shutDownTimeout time.Duration) *Server {
	return &Server{
		Port:            port,
		ShutDownTimeout: shutDownTimeout,
	}
}

// Run starts the HTTP server
func (s *Server) Run() {

	container := restful.DefaultContainer
	restful.Filter(globalLogging)
	services.Register(container)

	// TODO, docs can be enhanced using PostBuildSwaggerObjectHandler
	config := restfulspec.Config{
		WebServices: container.RegisteredWebServices(),
		APIPath:     "/apidocs.json"}

	container.Add(restfulspec.NewOpenAPIService(config))
	srv := &http.Server{Addr: fmt.Sprintf(":%d", s.Port)}

	allClosed := make(chan os.Signal, 1)

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
		<-sigterm
		log.Info("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), s.ShutDownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error(err, "error shutting down server")
		} else {
			log.Info("server stopped")
		}
		close(allClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err, "error at HTTP server")
	}
	<-allClosed
}

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	begin := time.Now()
	chain.ProcessFilter(req, resp)

	ms := float64(time.Now().Sub(begin).Nanoseconds() / 1000000.0)
	log.Info("",
		"elapsed", ms,
		"remote", req.Request.RemoteAddr,
		"uri", req.Request.RequestURI,
		"method", req.Request.Method,
		"code", resp.StatusCode(),
		"bytes", resp.ContentLength())
}
