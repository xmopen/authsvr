package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xmopen/authsvr/internal/config"
	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/authsvr/internal/endpoint"
	"github.com/xmopen/authsvr/internal/server"
	"github.com/xmopen/golib/pkg/xgoroutine"
	rpcserver "github.com/xmopen/gorpc/pkg/server"
)

type app struct {
	engine *gin.Engine
	apiSvr *http.Server
	rpcsvr *rpcserver.Server
	cancel context.CancelFunc
	close  chan error
	xlog   *xlogging.Entry
}

// init 初始化svr.
func (a *app) init(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		select {
		case r := <-sigs:
			a.close <- fmt.Errorf("syscall:[%+v]\n", r)
		}
	}()
	a.rpcsvr = rpcserver.NewServer()
	endpoint.Init(a.engine)
	server.Init(ctx, a.rpcsvr)
	a.run(ctx)
}

// run 运行svr.
func (a *app) run(ctx context.Context) {
	xgoroutine.SafeGoroutine(func() {
		network := config.Config().GetString("server.authsvr.rpc.network")
		addr := config.Config().GetString("server.authsvr.rpc.addr")
		a.xlog.Infof("rpc server running addr:[%+v] network:[%+v]", addr, network)
		if err := a.rpcsvr.Server(network, addr); err != nil {
			fmt.Printf("rpc server err:[%+v]\n", err)
		}
	})
	if err := a.apiSvr.ListenAndServe(); err != nil {
		a.close <- err
	}
}

func (a *app) quit() {
	select {
	case err := <-a.close:
		fmt.Println("svr done because err:" + err.Error())
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	r := gin.New()
	addr := config.Config().GetString("server.authsvr.http.addr")
	app := &app{
		engine: r,
		apiSvr: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
		cancel: cancel,
		close:  make(chan error, 1),
		xlog:   xlogging.Tag("authsvr.main"),
	}
	app.xlog.Infof("http server running addr:[%+v]", addr)
	app.init(ctx)
	app.quit()
}
