package brc20

import (
	"btc_order/internal/logic/brc20"
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route get /v1/market/latest brc20 GetTickerMarketInfo
//
// 获取指定Token的市场信息
//
// 获取指定Token的市场信息
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: TickerMarketReq
//
// Responses:
//  200: TickerMarketDetailResp

func GetTickerMarketInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TickerMarketReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := brc20.NewGetTickerMarketInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetTickerMarketInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
