package market

import (
	"btc_order/internal/errcode"
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"btc_order/pkg/apiclient"
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBtcPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBtcPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBtcPriceLogic {
	return &GetBtcPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type cmcResp struct {
	Data   map[string]types.CoinMarketInfo `json:"data"`
	Status types.CmcStatus                 `json:"status"`
}

func (l *GetBtcPriceLogic) GetBtcPrice() (resp *types.BtcCoinInfoResp, err error) {
	var (
		apiRes cmcResp
		apiErr types.CmcError
	)

	r := l.svcCtx.CmcApi.NewRequest(
		func(r *apiclient.Request) {
			r.SetQueryParam("id", "1")
			r.SetQueryParam("convert", "USD")
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

	data, ok := apiRes.Data["1"]
	if !ok {
		return nil, errcode.ErrCmcNotFound
	}

	resp = &types.BtcCoinInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: types.BtcCoinInfo{
			Price: data.Quote["USD"].Price,
		},
	}

	return
}
