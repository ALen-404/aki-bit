package market

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/market"
	"btc_order/internal/svc"
)

// swagger:route get /v1/market/btc/price market GetBtcPrice
//
// 获取BTC的当前价格
//
// 获取BTC的当前价格
//
// Responses:
//  200: BtcCoinInfoResp

func GetBtcPriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := market.NewGetBtcPriceLogic(r.Context(), svcCtx)
		resp, err := l.GetBtcPrice()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
