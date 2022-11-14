package api

import (
	"flag"
	"net/http"
	"time"

	"github.com/Pauca-Technologies/payment-ops-poc/bank-tranactional-ms/config"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
)

// Server server
type Server struct {
	Router *gin.Engine
	config *config.Config
}

// Controller controller
type Controller interface {
	SetupRoutes(router *gin.Engine)
}

// NewServer creates a new server instance
func NewServer(config *config.Config) *Server {
	server := &Server{
		Router: gin.Default(),
		config: config,
	}
	server.ConfigureLogging()
	return server
}

func (server *Server) ConfigureController(controller Controller) {
	controller.SetupRoutes(server.Router)
}

// ConfigureLogging configure gin logs
func (server *Server) ConfigureLogging() {
	flag.Parse()
	server.Router.Use(ginglog.Logger(3 * time.Second))
	server.Router.Use(gin.Recovery())
}

// Start runs the HTTP Server on a specific address
func (server *Server) Start() error {
	var readTimeout time.Duration
	readTimeout = time.Duration(server.config.ReadTimeout) * time.Second
	var writeTimeout time.Duration
	writeTimeout = time.Duration(server.config.WriteTimeout) * time.Second

	s := &http.Server{
		Addr:         server.config.ServerAddress,
		Handler:      server.Router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return s.ListenAndServe()
}
