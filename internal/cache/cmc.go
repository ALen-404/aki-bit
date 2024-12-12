package cache

import (
	"btc_order/internal/types"
	"btc_order/pkg/apiclient"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Cmc struct {
	api  *apiclient.Client
	rds  *redis.Redis
	quit chan struct{}
}

func NewCmc(rds *redis.Redis, api *apiclient.Client) *Cmc {
	return &Cmc{
		api:  api,
		rds:  rds,
		quit: make(chan struct{}),
	}
}

func (c *Cmc) Run() <-chan struct{} {
	init := make(chan struct{})

	go c.worker(init)

	return init
}

func (c *Cmc) Close() error {
	close(c.quit)

	return nil
}

type cmcListResp struct {
	Data   types.CategoryMarketData `json:"data"`
	Status types.CmcStatus          `json:"status"`
}

func (c *Cmc) worker(ch chan<- struct{}) {
	t := time.NewTicker(1 * time.Hour)
	started := false

	for {
		if started {
			select {
			case <-c.quit:
				return
			case <-t.C:
			}
		}

		var (
			apiRes cmcListResp
			apiErr types.CmcError
		)

		r := c.api.NewRequest(
			func(r *apiclient.Request) {
				r.SetQueryParam("id", apiclient.Brc20CategoryID)
				r.SetQueryParam("limit", "1")
				r.ForceContentType(apiclient.JsonContentType)
				r.SetResult(&apiRes)
				r.SetError(&apiErr)
			},
		)

		err := r.Execute(http.MethodGet, "/v1/cryptocurrency/category")
		if err != nil {
			logx.Errorw("request cmc occur unexpect error", logx.Field("error", err))
			time.Sleep(10 * time.Second)
			continue
		}

		if apiErr.Status.ErrCode != 0 {
			logx.Errorw("cmc service report an error",
				logx.Field("err_code", apiErr.Status.ErrCode),
				logx.Field("err_msg", apiErr.Status.ErrMessage),
			)
			time.Sleep(3 * time.Second)
			continue
		}

		totalNum := apiRes.Data.TokenNum
		cnt := 0

		var keySuffix string

		for {
			var (
				apiRes cmcListResp
				apiErr types.CmcError
			)

			r := c.api.NewRequest(
				func(r *apiclient.Request) {
					r.SetQueryParam("id", apiclient.Brc20CategoryID)
					r.SetQueryParam("start", strconv.Itoa(cnt*200+1))
					r.SetQueryParam("limit", "200")
					r.SetQueryParam("convert", "BTC")
					r.ForceContentType(apiclient.JsonContentType)
					r.SetResult(&apiRes)
					r.SetError(&apiErr)
				},
			)

			err := r.Execute(http.MethodGet, "/v1/cryptocurrency/category")
			if err != nil {
				logx.Errorw("request cmc occur unexpect error", logx.Field("error", err))
				time.Sleep(10 * time.Second)
				continue
			}

			if apiErr.Status.ErrCode != 0 {
				logx.Errorw("cmc service report an error",
					logx.Field("err_code", apiErr.Status.ErrCode),
					logx.Field("err_msg", apiErr.Status.ErrMessage),
				)
				time.Sleep(3 * time.Second)
				continue
			}

			// 以时间戳生成新的key
			keySuffix = strconv.FormatInt(time.Now().Unix(), 10)

			keyMap := make(map[string]struct{}, 0)
			keyAppend := func(keyPrefix string, score float64, val string) {
				key := fmt.Sprintf("%s%s", keyPrefix, keySuffix)

				_, err := c.rds.ZaddFloat(key, score, val)
				if err != nil {
					logx.Errorw("failed to update cmc market cache")
				}

				keyMap[key] = struct{}{}
			}

			for _, v := range apiRes.Data.Coins {
				data, err := json.Marshal(v)
				if err != nil {
					logx.Errorw("failed to marshal coin data", logx.Field("error", err))
					continue
				}

				if v.Quote["BTC"].MarketCap-0 < 1e-6 {
					keyAppend(CmcMarketCapPrefix, v.Quote["BTC"].FullyMarketCap, string(data))
				} else {
					keyAppend(CmcMarketCapPrefix, v.Quote["BTC"].MarketCap, string(data))
				}

				keyAppend(CmcVolume24hPrefix, v.Quote["BTC"].Volume24h, string(data))
				keyAppend(CmcPricePrefix, v.Quote["BTC"].Price, string(data))
			}

			keys := make([]any, 0, len(keyMap))
			for key := range keyMap {
				keys = append(keys, key)
			}

			_, err = c.rds.Rpush(CmcKeysKey, keys...)
			if err != nil {
				logx.Errorw("failed to update cmc cache", logx.Field("error", err))
				time.Sleep(5 * time.Second)
				continue
			}

			oldCnt, _ := c.rds.Llen(CmcKeysKey)
			if oldCnt > len(keyMap) {
				oldKeys, _ := c.rds.LpopCount(CmcKeysKey, len(keyMap))
				for _, v := range oldKeys {
					if _, ok := keyMap[v]; !ok {
						c.rds.Del(v)
					}
				}
			}

			cnt += 200
			if cnt >= totalNum {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}

		if !started {
			started = true
			close(ch)
		}
	}
}

const (
	CmcMarketCapPrefix = "cmc_marketcap_"
	CmcVolume24hPrefix = "cmc_volume24h_"
	CmcPricePrefix     = "cmc_price_"

	CmcKeysKey = "cmc_cached"
)
