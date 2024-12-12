//	btc_order
//
//	Description: btc_order service
//
//	Schemes: http, https
//	Host: localhost:3000
//	BasePath: /
//	Version: 0.0.1
//	SecurityDefinitions:
//	  Token:
//	    type: apiKey
//	    name: Authorization
//	    in: header
//	Security:
//	    - Token: []
//	Consumes:
//	  - application/json
//
//	Produces:
//	  - application/json
//
// swagger:meta
package main

import (
	"btc_order/internal/cache"
	"btc_order/internal/config"
	"btc_order/internal/cron"
	ierrcode "btc_order/internal/errcode"
	"btc_order/internal/handler"
	"btc_order/internal/svc"
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/suyuan32/simple-admin-common/utils/errcode"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/btc_order.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.CROSConf.Address))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	cmcCache := cache.NewCmc(ctx.Redis, ctx.CmcApi)
	cmcCache.Run()
	defer cmcCache.Close()

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		type HttpError struct {
			Message string `json:"msg"`
		}

		if ctx.Err() != nil {
			return http.StatusRequestTimeout, &HttpError{Message: "Request timeout"}
		}

		if errcode.IsGrpcError(err) {
			// don't unwrap error and get status.Message(),
			// it hides the rpc error headers.
			return errcode.CodeFromGrpcError(err), &HttpError{Message: err.Error()}
		}

		return ierrcode.CodeFromError(err), &HttpError{Message: err.Error()}
	})

	handler.RegisterHandlers(server, ctx)

	job := cron.New(ctx.PgDB)
	job.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
