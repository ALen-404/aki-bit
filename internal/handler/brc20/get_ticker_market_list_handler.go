package brc20

import (
	"btc_order/internal/logic/brc20"
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route get /v1/market/list/latest brc20 GetTickerMarketList
//
// 获取已知Token列表的市场信息
//
// 获取已知Token列表的市场信息
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: TickerMarketListReq
//
// Responses:
//  200: TickerMarketListResp

func GetTickerMarketListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TickerMarketListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := brc20.NewGetTickerMarketListLogic(r.Context(), svcCtx)
		resp, err := l.GetTickerMarketList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
