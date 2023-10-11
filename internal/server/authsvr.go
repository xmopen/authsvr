package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/xmopen/authsvr/internal/models/usermodels"

	"github.com/xmopen/authsvr/internal/service/authservice"

	"github.com/xmopen/commonlib/pkg/server/authserver"
)

var (
	authSvrInstance authserver.IAuthServer
	authSvrOnce     sync.Once
)

// AuthSvr auth svr
type AuthSvr struct {
}

// NewAuthSvr return a authsvr instance
func NewAuthSvr() authserver.IAuthServer {
	authSvrOnce.Do(func() {
		authSvrInstance = &AuthSvr{}
	})
	return authSvrInstance
}

// GetUserInfoByAccount get user info by account
func (a *AuthSvr) GetUserInfoByAccount(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	if request.XMAccount == "" {
		return fmt.Errorf("request xmaccount is empty")
	}
	xmUser, err := authservice.Service().XMUserWithAccount(request.XMAccount)
	if err != nil {
		return err
	}
	if xmUser == nil {
		xmUser, err = usermodels.XMUserWithAccount(request.XMAccount)
		if err != nil || xmUser == nil {
			return fmt.Errorf("user empty account:[%+v] err:[%+v]", request.XMAccount, err)
		}
		xgoroutine.SafeGoroutine(func() {
			if err = authservice.Service().CreateXMUserToRedis(xmUser); err != nil {
				xlogging.Tag("authsvr.gorpc.server.loc").Errorf("create xmuser to redis err:[%+v] user:[%+v]",
					err, xmUser)
			}
		})
	}
	response.XMUserInfo = xmUser
	return nil
}

// GetUserInfoByToken get user info by token
func (a *AuthSvr) GetUserInfoByToken(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	if request.XMToken == "" {
		return fmt.Errorf("xmtoken is empty")
	}
	xmUser, err := authservice.Service().XMUserWithToken(request.XMToken)
	if err != nil {
		return err
	}
	if xmUser == nil {
		return fmt.Errorf("xm user is nil token:[%+v]", request.XMToken)
	}
	response.XMUserInfo = xmUser
	return nil
}
