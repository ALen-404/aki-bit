package price

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/btc/price"
	"btc_order/internal/svc"
)

// swagger:route get /v1/btc/price btc_price GetBtcPirce
//
// 获取btc价格
//
// 获取btc价格
//
// Responses:
//  200: BtcPriceResp

func GetBtcPirceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := price.NewGetBtcPirceLogic(r.Context(), svcCtx)
		resp, err := l.GetBtcPirce()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
