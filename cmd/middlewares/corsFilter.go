package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
}

func CorsFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		if origin != "" && allowedOrigins[origin] == true {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")

			// If you return a specific allowed origin in the Access-Control-Allow-
			// Origin response header, then you should also include a Vary: Origin header
			// to ensure the browser and any network proxies only cache the response for
			// this specific requesting origin.
			ctx.Header("Vary", "Origin")
		}

		// if it's a preflight request
		if isPreFlightRequest(ctx) {
			// if origin is not allowd, then reject the preflight request
			if origin == "" || !allowedOrigins[origin] {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}

			// if origin is allowed, then allow the preflight request
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, DELETE")

			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func isPreFlightRequest(ctx *gin.Context) bool {
	return ctx.Request.Method == "OPTIONS" && ctx.GetHeader("Access-Control-Request-Method") != ""
}
