package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/authsvr/internal/endpoint/middleware"
	"github.com/xmopen/authsvr/internal/endpoint/userauth"
)

// Init  gin router.
func Init(r *gin.Engine) {
	middleware.Init(r)

	authAPI := userauth.New()
	group := r.Group("/openxm/api/v1/auth")
	group.POST("/login", authAPI.UserLogin)
	group.POST("/register", authAPI.UserRegisterAndLogin)
}
