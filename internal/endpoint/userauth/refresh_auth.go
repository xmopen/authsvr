package userauth

import (
	"net/http"
	"sync"

	"github.com/xmopen/commonlib/pkg/database/xmuser"

	"github.com/xmopen/authsvr/internal/service/authservice"
	"github.com/xmopen/authsvr/internal/util/apputils"

	"github.com/xmopen/commonlib/pkg/errcode"

	"github.com/gin-gonic/gin"
)

var (
	refreshAPIInstance *RefreshAuthAPI
	refreshAPIOnce     sync.Once
)

type userRefreshResponse struct {
	XMToken  string         `json:"xm_token"`
	UserInfo *xmuser.XMUser `json:"user_info"`
}

// RefreshAuthAPI refresh auth api
type RefreshAuthAPI struct {
}

// NewRefreshAuthAPI return a single refresh api instance
func NewRefreshAuthAPI() *RefreshAuthAPI {
	refreshAPIOnce.Do(func() {
		refreshAPIInstance = &RefreshAuthAPI{}
	})
	return refreshAPIInstance
}

// CheckXMUserWithToken check xmUser with token
// and refreshing a token when it is about to expire
func (r *RefreshAuthAPI) CheckXMUserWithToken(c *gin.Context) {
	xmToken := c.Query("xm_token")
	if xmToken == "" {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog := apputils.Log(c)
	xmUser, err := authservice.Service().XMUserWithToken(xmToken)
	if err != nil {
		xlog.Errorf("get xmuser with token err:[%+v] token:[%+v]", err, xmToken)
		// TODO: replace errcode
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	if xmUser == nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	newToken := authservice.Service().RefreshXMToken(xlog, xmToken, xmUser)
	if newToken != "" {
		xmToken = newToken
	}
	c.JSON(http.StatusOK, errcode.Success(&userRefreshResponse{
		UserInfo: xmUser,
		XMToken:  xmToken,
	}))
}
