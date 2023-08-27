package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/golib/pkg/middleware"
)

// Init gin middleware.
func Init(r *gin.Engine) {
	r.Use(middleware.Cors())
}
