package middleware

import (
	"Lumino/store"
	"github.com/gin-gonic/gin"
)

// DB is a middleware function that initializes the Datastore and attaches to the context of every http.Request.
func DB(v *store.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
