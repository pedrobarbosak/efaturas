package server

import (
	"fmt"
	"time"

	"efaturas-xtreme/internal/api"
	"efaturas-xtreme/internal/api/middlewares"
	"efaturas-xtreme/internal/auth"
	"efaturas-xtreme/internal/service"
	"efaturas-xtreme/internal/service/repository"
	"efaturas-xtreme/pkg/db"
	"efaturas-xtreme/pkg/efaturas"
	"efaturas-xtreme/pkg/sse"

	"efaturas-xtreme/pkg/errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Run() error
}

type server struct {
	cfg Config

	db       *db.DB
	sse      sse.Server
	efaturas efaturas.Service
	auth     *auth.Service

	service service.Service
	api     *api.API
}

func (s *server) Run() error {
	r := gin.New()
	r.Use(gzip.Gzip(gzip.BestCompression)) // needs to go before recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-type", "authorization", "Content-Length", "Content-Language", "TE", "Content-Disposition", "User-Agent", "Referrer", "Host", "Access-Control-Allow-Origin", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Accept-Encoding", "Accept-Language", "X-Shared", ""},
		ExposeHeaders:    []string{"Authorization", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(location.Default())
	r.Use(requestid.New())
	r.Use(secure.New(secure.Config{
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
		FrameDeny:          true,
	}))
	r.Use(errors.Middleware())

	apiGroup := r.Group(s.cfg.API.Version)
	authMiddleware := middlewares.Auth(s.auth)

	s.api.Init(apiGroup, authMiddleware)

	addr := fmt.Sprintf(":%d", s.cfg.API.Port)
	return r.Run(addr)
	// return http.ListenAndServeTLS(addr, path.Join(s.cfg.API.Certificates, "server.crt"), path.Join(s.cfg.API.Certificates, "server.key"), r)
}

func (s *server) setupDatabase() error {
	db, err := db.New(s.cfg.Database.URI, s.cfg.Database.Name)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func New(cfg Config) (Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.New("invalid config:", err)
	}

	s := &server{cfg: cfg}

	if err := s.setupDatabase(); err != nil {
		return nil, err
	}

	s.sse = sse.New()
	s.efaturas = efaturas.New()
	s.auth = auth.New(s.efaturas)

	repo, err := repository.New(s.db)
	if err != nil {
		return nil, err
	}

	s.service = service.New(repo, s.efaturas, s.sse)
	s.api = api.New(s.service, s.auth, s.sse)

	return s, nil
}
