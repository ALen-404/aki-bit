package brc20

import (
	"btc_order/internal/errcode"
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"btc_order/pkg/apiclient"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTickerMarketInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTickerMarketInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTickerMarketInfoLogic {
	return &GetTickerMarketInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type cmcQuoteResp struct {
	Data   map[string]types.CoinMarketInfo `json:"data"`
	Status types.CmcStatus                 `json:"status"`
}

type cmcOhlcvResp struct {
	Data   types.OhlcvSeriesData `json:"data"`
	Status types.CmcStatus       `json:"status"`
}

func (l *GetTickerMarketInfoLogic) GetTickerMarketInfo(req *types.TickerMarketReq) (resp *types.TickerMarketDetailResp, err error) {
	var (
		apiRes cmcQuoteResp
		apiErr types.CmcError
	)

	r := l.svcCtx.CmcApi.NewRequest(
		func(r *apiclient.Request) {
			r.SetQueryParam("id", req.ID)
			r.SetQueryParam("convert", "BTC")
			r.ForceContentType(apiclient.JsonContentType)
			r.SetResult(&apiRes)
			r.SetError(&apiErr)
		},
	)

	err = r.Execute(http.MethodGet, "/v2/cryptocurrency/quotes/latest")
	if err != nil {
		logx.Errorw("failed to request cmc coin quote", logx.Field("error", err))

		switch {
		case apiErr.Status.ErrCode != 0:
			return nil, errcode.ErrCmcApiError
		case errors.As(err, &apiclient.ApiError{}):
			return nil, errcode.ErrHttpError
		default:
			return nil, err
		}
	}

	data, ok := apiRes.Data[req.ID]
	if !ok {
		return nil, errcode.ErrCmcNotFound
	}

	resp = &types.TickerMarketDetailResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
	}

	resp.Data.ID = data.ID
	resp.Data.Tick = data.Symbol
	resp.Data.MarketCap = data.Quote["BTC"].MarketCap
	resp.Data.Volume24h = data.Quote["BTC"].Volume24h
	resp.Data.SatPrice = data.Quote["BTC"].Price

	// 请求k线图
	var ohlcvResp cmcOhlcvResp
	r = l.svcCtx.CmcApi.NewRequest(
		func(r *apiclient.Request) {
			r.SetQueryParam("id", req.ID)
			r.SetQueryParam("convert", "BTC")
			r.SetQueryParam("count", strconv.Itoa(req.Count))
			r.SetQueryParam("interval", req.Interval)
			r.ForceContentType(apiclient.JsonContentType)
			r.SetResult(&ohlcvResp)
			r.SetError(&apiErr)
		},
	)

	err = r.Execute(http.MethodGet, "/v2/cryptocurrency/ohlcv/historical")
	if err != nil {
		var retErr error

		switch {
		case apiErr.Status.ErrCode != 0:
			// 如果key不支持，跳过
			if apiErr.Status.ErrCode == 1006 {
				err = nil
				return
			}
			retErr = errcode.ErrCmcApiError
		case errors.As(err, &apiclient.ApiError{}):
			retErr = errcode.ErrHttpError
		default:
			retErr = err
		}

		if retErr != nil {
			logx.Errorw("failed to request cmc coin quote", logx.Field("error", retErr))

			return nil, retErr
		}
	}

	return
}
