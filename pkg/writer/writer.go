package writer

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r ResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func NewResponseWriter(writer gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		body:           &bytes.Buffer{},
		ResponseWriter: writer,
	}
}

// WriteResponse writes a json response with api required headers
func WriteResponse(c *gin.Context, status int, value interface{}) {
	c.Header("Content-Type", "application/json")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "0")
	c.Header("Cache-Control", "no-store")

	c.JSON(status, value)
	return
}

// WriteError writes an error response
func WriteError(c *gin.Context, status int, err error) {
	if status == http.StatusInternalServerError {
		WriteResponse(c, status, map[string]string{"error": "internal server error"})
	}
	WriteResponse(c, status, map[string]string{"error": err.Error()})
}
