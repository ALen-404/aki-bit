package svc

import (
	"btc_order/internal/config"
	"btc_order/internal/middleware"
	"btc_order/pkg/apiclient"
	"btc_order/pkg/postgresql"
	"net/http"
	"time"

	cache "github.com/victorspringer/http-cache"
	redisAdapter "github.com/victorspringer/http-cache/adapter/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config config.Config
	PgDB   *postgresql.Database
	Redis  *redis.Redis
	CmcApi *apiclient.Client

	Cache     rest.Middleware
	RateLimit rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	pgdb := postgresql.NewPostgresqlDriver(&c.PostgresqlDB)

	rds := redis.MustNewRedis(c.RedisConf)

	apic := apiclient.New(
		func(ac *apiclient.Client) {
			if c.CmcApi.UseHttps {
				ac.Schema = "https://"
			} else {
				ac.Schema = "http://"
			}
			ac.Domain = c.CmcApi.Domain
			ac.Header = &http.Header{
				"X-CMC_PRO_API_KEY": {c.CmcApi.Key},
			}
			ac.Retry = 3
			ac.RetryTime = 10000
		},
	)

	ringOpt := &redisAdapter.RingOptions{
		Addrs: map[string]string{
			"server": rds.Addr,
		},
	}
	httpCache, err := cache.NewClient(
		cache.ClientWithAdapter(redisAdapter.NewAdapter(ringOpt)),
		cache.ClientWithTTL(3*time.Minute),
		cache.ClientWithRefreshKey("pvreset"),
	)
	logx.Must(err)

	return &ServiceContext{
		Config:    c,
		PgDB:      pgdb,
		Redis:     rds,
		CmcApi:    apic,
		Cache:     middleware.NewCacheMiddleware(httpCache).Handle,
		RateLimit: middleware.NewRateLimitMiddleware().Handle,
	}
}
