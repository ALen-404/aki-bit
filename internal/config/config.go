package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	RedisConf    redis.RedisConf
	CROSConf     CROSConf
	PostgresqlDB PgDB
	CmcApi       CmcApiConf
}

type CROSConf struct {
	Address string `json:",env=CROS_ADDRESS"`
}

type PgDB struct {
	Addr     string `json:"Addr"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

type CmcApiConf struct {
	Key      string
	Domain   string `json:",default=pro-api.coinmarketcap.com"`
	UseHttps bool   `json:",default=true"`
}
