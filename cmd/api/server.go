package api

import (
	"net/http"

	"github.com/Luckny/space-it/cmd/middlewares"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("accesslvl", middlewares.ValidAccessLvl)
	}

	ginDefault := gin.Default()
	router := ginDefault.Group("/api/v1")

	router.Use(middlewares.EnsureJSONContentType())
	router.Use(middlewares.RateGuard(server.Limiter))
	router.Use(middlewares.Authenticate(store))
	router.Use(middlewares.AuditLogger(store))

	router.POST("/users", server.registerUser)

	router.Use(middlewares.RequireAuthentication())
	router.POST("/spaces", server.createSpace)

	router.POST(
		"/spaces/:spaceID/messages",
		middlewares.RequireAccessLvl(middlewares.WriteAccess, store),
		func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, "not implemented")
		},
	)

	router.POST(
		"/spaces/:spaceID/members",
		middlewares.RequireAccessLvl(middlewares.AdminAccess, store),
		server.addMemberToSpace,
	)

	server.Router = ginDefault
	return server
}

// Run's the api server
func (server *Server) Run(addr string) error {
	return server.Router.Run(addr)
}
