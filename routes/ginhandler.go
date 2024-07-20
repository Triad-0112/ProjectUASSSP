package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinHandlerToHTTP(handler func(*gin.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new Gin context from the request
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		handler(c)
	}
}
