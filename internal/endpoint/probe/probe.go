package probe

import (
	"net/http"

	"github.com/xmopen/commonlib/pkg/apphelper/ginhelper"

	"github.com/gin-gonic/gin"
)

// HealthProbe Kubernetes Health Probe.
func HealthProbe(c *gin.Context) {
	ginhelper.Log(c).Infof("authsvr health probe")
	c.JSON(http.StatusOK, "success")
}
