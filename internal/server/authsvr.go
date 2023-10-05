package server

import (
	"context"
	"sync"

	"github.com/xmopen/commonlib/pkg/database/xmuser"

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

func (a *AuthSvr) GetUserInfoByAccount(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	response.XMUserInfo = &xmuser.XMUser{
		UserName: "zhenxinma",
	}
	return nil
}

func (a *AuthSvr) GetUserInfoByToken(ctx context.Context, request *authserver.AuthSvrRequest, response *authserver.AuthSvrResponse) error {
	response.XMUserInfo = &xmuser.XMUser{
		UserName: "zhenxinma",
	}
	return nil
}
