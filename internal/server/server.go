package server

import (
	"context"

	"github.com/xmopen/commonlib/pkg/server/authserver"
	"github.com/xmopen/golib/pkg/xlogging"
	rpcserver "github.com/xmopen/gorpc/pkg/server"
)

// Init rpcsvr init
func Init(ctx context.Context, rpcsvr *rpcserver.Server) {
	rpcsvr.SetTrace(true)
	rpcsvr.RegisterName(authserver.AuthSvrName, NewAuthSvr(), "")
	xlogging.Tag("rpcsvr").Infof("init success")
}
