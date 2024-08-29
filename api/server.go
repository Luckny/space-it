package api

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Server struct {
	store   db.Store
	router  *gin.Engine
	limiter *rate.Limiter
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store, limiter: rate.NewLimiter(rate.Limit(2), 2)}
	ginDefault := gin.Default()
	router := ginDefault.Group("/api/v1")

	router.Use(EnsureJSONContentType())
	router.Use(RateGuard(server.limiter))
	router.Use(server.Authenticate())

	router.POST("/users", server.registerUser)

	server.router = ginDefault
	return server
}

// Run's the api server
func (server *Server) Run(addr string) error {
	return server.router.Run(addr)
}

// writeResponse writes a json response with api required headers
func writeResponse(c *gin.Context, status int, value interface{}) {
	c.Header("Content-Type", "application/json")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "0")
	c.Header("Cache-Control", "no-store")

	c.JSON(status, value)
	return
}

// writeError writes an error response
func writeError(c *gin.Context, status int, err error) {
	if status == http.StatusInternalServerError {
		writeResponse(c, status, map[string]string{"error": "internal server error"})
	}
	writeResponse(c, status, map[string]string{"error": err.Error()})
}
