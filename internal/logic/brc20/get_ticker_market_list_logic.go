package brc20

import (
	"btc_order/internal/cache"
	"btc_order/internal/errcode"
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"btc_order/pkg/apiclient"
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTickerMarketListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTickerMarketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTickerMarketListLogic {
	return &GetTickerMarketListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type cmcListResp struct {
	Data   types.CategoryMarketData `json:"data"`
	Status types.CmcStatus          `json:"status"`
}

func (l *GetTickerMarketListLogic) GetTickerMarketList(req *types.TickerMarketListReq) (resp *types.TickerMarketListResp, err error) {
	keys, err := l.svcCtx.Redis.Lrange(cache.CmcKeysKey, 0, 2)
	if err != nil {
		logx.Errorw("failed to read cmc cache", logx.Field("error", err))
		return nil, errcode.ErrCmcNotFound
	}

	if req.Page == 0 {
		req.Page = 1
	}

	keyFind := func(prefix string) (string, bool) {
		for _, v := range keys {
			if strings.HasPrefix(v, prefix) {
				return v, true
			}
		}

		return "", false
	}

	resp = &types.TickerMarketListResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
	}

	var (
		key   string
		found bool
	)

	switch strings.ToLower(req.Sort) {
	case "marketcap":
		key, found = keyFind(cache.CmcMarketCapPrefix)
	case "volume24h":
		key, found = keyFind(cache.CmcVolume24hPrefix)
	case "price":
		key, found = keyFind(cache.CmcPricePrefix)
	default:
		var (
			apiRes cmcListResp
			apiErr types.CmcError
		)

		r := l.svcCtx.CmcApi.NewRequest(
			func(r *apiclient.Request) {
				r.SetQueryParam("id", apiclient.Brc20CategoryID)
				r.SetQueryParam("limit", "200")
				r.SetQueryParam("convert", "BTC")
				r.ForceContentType(apiclient.JsonContentType)
				r.SetResult(&apiRes)
				r.SetError(&apiErr)
			},
		)

		err = r.Execute(http.MethodGet, "/v1/cryptocurrency/category")
		if err != nil {
			logx.Errorw("failed to request cmc category info", logx.Field("error", err))

			switch {
			case apiErr.Status.ErrCode != 0:
				return nil, errcode.ErrCmcApiError
			case errors.As(err, &apiclient.ApiError{}):
				return nil, errcode.ErrHttpError
			default:
				return nil, err
			}
		}

		resp.Data.Total = uint64(apiRes.Data.TokenNum)
		for _, coin := range apiRes.Data.Coins {
			marketCap := coin.Quote["BTC"].MarketCap

			if coin.Quote["BTC"].MarketCap-0 < 0.01 {
				marketCap = coin.Quote["BTC"].FullyMarketCap
			}

			resp.Data.Data = append(resp.Data.Data, types.TickerMarketInfo{
				ID:        coin.ID,
				Tick:      coin.Symbol,
				MarketCap: marketCap,
				Volume24h: coin.Quote["BTC"].Volume24h,
				SatPrice:  coin.Quote["BTC"].Price,
			})
		}

		return
	}

	if !found {
		return nil, errcode.ErrCmcNotFound
	}

	total, err := l.svcCtx.Redis.Zcount(key, math.MinInt64, math.MaxInt64)
	if err != nil {
		logx.Errorw("failed to read cmc cache", logx.Field("error", err))
		return nil, errcode.ErrCmcNotFound
	}

	resp.Data.Total = uint64(total)

	data, err := l.svcCtx.Redis.ZrevrangebyscoreWithScoresByFloatAndLimit(key, 0, math.MaxFloat64, int(req.Page-1), int(req.PageSize))
	if err != nil {
		logx.Errorw("failed to read cmc cache", logx.Field("error", err))
		return nil, errcode.ErrCmcNotFound
	}

	for _, p := range data {
		var coin types.CoinMarketInfo

		err := json.Unmarshal([]byte(p.Key), &coin)
		if err != nil {
			logx.Errorw("failed to unmarshal cached market info", logx.Field("error", err))
			continue
		}

		marketCap := coin.Quote["BTC"].MarketCap

		if coin.Quote["BTC"].MarketCap-0 < 0.01 {
			marketCap = coin.Quote["BTC"].FullyMarketCap
		}

		resp.Data.Data = append(resp.Data.Data, types.TickerMarketInfo{
			ID:        coin.ID,
			Tick:      coin.Symbol,
			MarketCap: marketCap,
			Volume24h: coin.Quote["BTC"].Volume24h,
			SatPrice:  coin.Quote["BTC"].Price,
		})
	}

	return
}
