package userauth

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/authsvr/internal/models/usermodels"
	"github.com/xmopen/authsvr/internal/service/authservice"
	"github.com/xmopen/authsvr/internal/util/apputils"
	"github.com/xmopen/commonlib/pkg/database/xmuser"
	"github.com/xmopen/commonlib/pkg/errcode"
	"github.com/xmopen/golib/pkg/utils/commonutil"
	"github.com/xmopen/golib/pkg/xlogging"
)

// defaultUserRegisterIcon 默认用户注册头像
const defaultUserRegisterIcon = "https://typoraimg-1303903194.cos.ap-guangzhou.myqcloud.com/%E5%A4%B4%E5%83%8F.jpg"

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
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type userLoginResponse struct {
	XMToken  string         `json:"xm_token"`
	UserInfo *xmuser.XMUser `json:"user_info"`
}

// UserLogin login,如果用户未注册则注册.
func (a *API) UserLogin(c *gin.Context) {
	request, err := a.unmarshalCheckLoginOrRegisterParam(c)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog := apputils.Log(c)
	loginResponse, err := a.loginWithAccount(xlog, request)
	if err != nil {
		xlog.Errorf("login fail err:[%+v]", err)
		c.JSON(http.StatusOK, errcode.ErrorUserLoginFail)
		return
	}
	// 密码不正确.
	if loginResponse == nil {
		c.JSON(http.StatusOK, errcode.ErrorUserLoginFail)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(loginResponse))
}

// UserRegisterAndLogin 用户注册并且登录.
func (a *API) UserRegisterAndLogin(c *gin.Context) {
	request, err := a.unmarshalCheckLoginOrRegisterParam(c)
	if err != nil {
		c.JSON(http.StatusOK, errcode.ErrorParam)
		return
	}
	xlog := apputils.Log(c)
	loginResponse, err := a.userRegistry(xlog, request)
	if err != nil {
		xlog.Errorf("user registry err:[%+v] request:[%+v]", err, request)
		c.JSON(http.StatusOK, errcode.ErrorUserRegisterFail)
		return
	}
	c.JSON(http.StatusOK, errcode.Success(loginResponse))
}

// unmarshalCheckLoginOrRegisterParam 序列化Request并且校验参数
func (a *API) unmarshalCheckLoginOrRegisterParam(c *gin.Context) (*userLoginRequest, error) {
	request := &userLoginRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		return nil, err
	}
	if request.Account == "" || request.Password == "" {
		return nil, fmt.Errorf("param is empty")
	}
	if len(request.Account) < 6 || len(request.Account) > 12 {
		return nil, fmt.Errorf("account is illegal")
	}
	match, err := regexp.Match(".{6,12}", []byte(request.Account))
	if err != nil || !match {
		return nil, fmt.Errorf("account is illegal")
	}
	return request, nil
}

func (a *API) loginWithAccount(xlog *xlogging.Entry, request *userLoginRequest) (*userLoginResponse, error) {
	xmUser, err := a.verifyUserIsRegistry(request)
	if err != nil {
		return nil, err
	}
	// 用户未注册.
	if xmUser == nil {
		return nil, nil
	}
	if xmUser.UserPassword != commonutil.MD5(request.Password) {
		return nil, nil
	}
	// 校验通过.
	xmToken, err := authservice.Service().CreateXMUserToken(xmUser)
	if err != nil {
		return nil, err
	}
	return &userLoginResponse{
		XMToken:  xmToken,
		UserInfo: xmUser,
	}, nil
}

// verifyUserIsRegistry 校验用户是否注册.
func (a *API) verifyUserIsRegistry(req *userLoginRequest) (*xmuser.XMUser, error) {
	// 1、查缓存中是否有用户信息.
	xmUser, err := authservice.Service().XMUserWithAccount(req.Account)
	if err != nil {
		return nil, err
	}
	if xmUser != nil {
		return xmUser, nil
	}
	// 2、查询DB中是否有用户信息.
	xmUser, err = usermodels.XMUserWithAccount(req.Account)
	if err != nil {
		return nil, err
	}
	return xmUser, nil
}

// userRegistry 注册用户.
// 密码规则: 6-12位引文字符.
func (a *API) userRegistry(xlog *xlogging.Entry, request *userLoginRequest) (*userLoginResponse, error) {
	xmUser, err := usermodels.XMUserWithAccount(request.Account)
	if err != nil {
		xlog.Errorf("get xmuser with account err:[%+v] request:[%+v]", err, request)
	}
	if xmUser != nil {
		return nil, fmt.Errorf("the account has been registerd")
	}
	xmUser = &xmuser.XMUser{
		UserAccount:  request.Account,
		UserPassword: commonutil.MD5(request.Password),
		UserName:     request.Name,
		UserIcon:     defaultUserRegisterIcon,
		CreateTime:   time.Now(),
		LastLogin:    time.Now(),
	}
	if err := usermodels.SaveUser(xmUser); err != nil {
		return nil, err
	}
	xmToken, err := authservice.Service().CreateXMUserToken(xmUser)
	if err != nil {
		return nil, err
	}
	// 同时登录.
	return &userLoginResponse{
		XMToken:  xmToken,
		UserInfo: xmUser,
	}, nil
}
