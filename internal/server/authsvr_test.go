package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/xmopen/commonlib/pkg/server/authserver"
	"github.com/xmopen/gorpc/pkg/client"
)

func TestAuthSvr(t *testing.T) {
	c, err := client.NewClient("tcp", ":18849", nil)
	if err != nil {
		panic(err)
	}
	c.Trace = true
	request := &authserver.AuthSvrRequest{}
	response := &authserver.AuthSvrResponse{}
	if err = c.Call(context.TODO(), authserver.AuthSvrName, string(authserver.AuthSvrMethodTypeOfGetUserInfoByXMAccount),
		request, response); err != nil {
		panic(err)
	}
	fmt.Printf("response:[%+v]\n", response)
}
