package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/xmopen/authsvr/internal/endpoint/middleware"
	"github.com/xmopen/authsvr/internal/endpoint/probe"
	"github.com/xmopen/authsvr/internal/endpoint/userauth"
)

// Init  gin router.
func Init(r *gin.Engine) {
	middleware.Init(r)

	router := r.Group("/openxm/api/v1/auth")
	router.GET("/probe", probe.HealthProbe)

	authAPI := userauth.New()
	router.POST("/login", authAPI.UserLogin)
	router.POST("/register", authAPI.UserRegisterAndLogin)

	refreshAPI := userauth.NewRefreshAuthAPI()
	router.GET("/check/session", refreshAPI.CheckXMUserWithToken)
}
