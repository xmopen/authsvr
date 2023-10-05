package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		if err := a.rpcsvr.Server("tcp", ":18849"); err != nil {
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
	app := &app{
		engine: r,
		apiSvr: &http.Server{
			Addr:              ":8849",
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
		cancel: cancel,
		close:  make(chan error, 1), // 容量为1不阻塞.
	}

	app.init(ctx)
	app.quit()
}
