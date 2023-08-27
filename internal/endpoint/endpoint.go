package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/authsvr/internal/endpoint/middleware"
)

// Init  gin router.
func Init(r *gin.Engine) {
	middleware.Init(r)
}
