package api

import (
	"github.com/Luckny/space-it/cmd/middlewares"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Server struct {
	store   db.Store
	Router  *gin.Engine
	Limiter *rate.Limiter
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:   store,
		Limiter: rate.NewLimiter(rate.Limit(2), 2),
	}
	ginDefault := gin.Default()
	router := ginDefault.Group("/api/v1")

	router.Use(middlewares.EnsureJSONContentType())
	router.Use(middlewares.RateGuard(server.Limiter))
	router.Use(middlewares.Authenticate(store))
	router.Use(middlewares.AuditLogger(store))

	router.POST("/users", server.registerUser)

	server.Router = ginDefault
	return server
}

// Run's the api server
func (server *Server) Run(addr string) error {
	return server.Router.Run(addr)
}
