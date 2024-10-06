package api

import (
	"net/http"

	"github.com/Luckny/space-it/cmd/middlewares"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/config"
	"github.com/Luckny/space-it/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/time/rate"
)

type Server struct {
	store      db.Store
	Router     *gin.Engine
	Limiter    *rate.Limiter
	tokenMaker token.Maker
	Config     config.Config
}

func NewServer(store db.Store, config config.Config) *Server {
	server := &Server{
		store:      store,
		Limiter:    rate.NewLimiter(rate.Limit(2), 2),
		tokenMaker: token.NewCookieStore(),
		Config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("accesslvl", middlewares.ValidAccessLvl)
	}

	router := gin.Default()

	router.Use(middlewares.EnsureJSONContentType())
	router.Use(middlewares.RateGuard(server.Limiter))

	// CORS preflight requests should
	// be handled before API requests authentication because credentials are never
	// sent on a preflight request, so it would always fail otherwise
	router.Use(middlewares.CorsFilter())

	// at least one of the following two middleware should succeed
	// for user to be authenticated
	router.Use(middlewares.Authenticate(store))
	router.Use(middlewares.VerifyToken(server.tokenMaker))

	// log all requests
	router.Use(middlewares.AuditLogger(store))

	router.POST(makeUrl("/users"), server.registerUser)

	// require that all following requests require authentication
	router.Use(middlewares.RequireAuthentication())

	router.POST(makeUrl("/users/login"), server.loginUser)
	router.DELETE(makeUrl("/users/logout"), server.logoutUser)
	router.POST(makeUrl("/spaces"), server.createSpace)

	router.GET(makeUrl("/test"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	router.POST(
		makeUrl("/spaces/:spaceID/messages"),
		middlewares.RequireAccessLvl(middlewares.WriteAccess, store),
		func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, "not implemented")
		},
	)

	router.POST(
		makeUrl("/spaces/:spaceID/members"),
		middlewares.RequireAccessLvl(middlewares.AdminAccess, store),
		server.addMemberToSpace,
	)

	server.Router = router
	return server
}

// Run's the api server
func (server *Server) Run(addr string) error {
	return server.Router.RunTLS(addr, "cert.pem", "key.pem")
}

func makeUrl(path string) string {
	return "/api/v1" + path
}
