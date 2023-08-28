package userauth

import (
	"net/http"

	"github.com/xmopen/authsvr/internal/models/usermodels"

	"github.com/xmopen/authsvr/internal/service/authservice"

	"github.com/xmopen/commonlib/pkg/database/xmuser"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/commonlib/pkg/errcode"
)

// API user auth api.
type API struct {
}

// New a api instance.
func New() *API {
	return &API{}
}

type userLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type userLoginResponse struct {
	XMToken  string         `json:"xm_token"`
	UserInfo *xmuser.XMUser `json:"user_info"`
}

// UserLoginWithRegister login,如果用户未注册则注册.
func (a *API) UserLoginWithRegister(c *gin.Context) {
	request := &userLoginRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	if request.Account == "" || request.Password == "" {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}

}

func (a *API) loginWithRegister(req *userLoginRequest) (*userLoginResponse, error) {
	isRegistry, err := a.verifyUserIsRegistry(req)
	if err != nil {
		return nil, err
	}
	if !isRegistry {
		// 用户未注册.
	}
	// 用户注册.
	return nil, nil
}

// verifyUserIsRegistry 校验用户是否注册.
func (a *API) verifyUserIsRegistry(req *userLoginRequest) (bool, error) {
	// 1、查缓存中是否有用户信息.
	xmUser, err := authservice.Service().XMUserWithAccount(req.Account)
	if err != nil {
		return false, err
	}
	if xmUser != nil {
		return true, nil
	}
	// 2、查询DB中是否有用户信息.
	xmUser, err = usermodels.XMUserWithAccount(req.Account)
	if err != nil {
		return false, err
	}
	return xmUser != nil, nil
}

// userRegistry 注册用户.
func (a *API) userRegistry(req *userLoginRequest) (*xmuser.XMUser, error) {

	// 1、校验账号是否符合规则.

	// 2、MD5加密密码.
	return nil, nil
}
