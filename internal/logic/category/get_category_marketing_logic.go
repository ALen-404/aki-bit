package category

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

type GetCategoryMarketingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryMarketingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryMarketingLogic {
	return &GetCategoryMarketingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type cmcResp struct {
	Data   types.CategoryMarketListData `json:"data"`
	Status types.CmcStatus              `json:"status"`
}

func (l *GetCategoryMarketingLogic) GetCategoryMarketing() (resp *types.CategoryMarketInfoResp, err error) {
	var (
		apiRes cmcResp
		apiErr types.CmcError
	)

	r := l.svcCtx.CmcApi.NewRequest(
		func(r *apiclient.Request) {
			r.SetQueryParam("id", apiclient.Brc20CategoryID)
			r.SetQueryParam("limit", "1")
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

	resp = &types.CategoryMarketInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: types.CategoryMarketInfo{
			MarketCap: apiRes.Data.MarketCap,
			Volume24h: apiRes.Data.Volume,
		},
	}

	return
}
