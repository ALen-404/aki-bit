package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/btc/order"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route post /v1/btc/order btc_order CreatedOrdOrder
//
// 创建订单
//
// 创建订单
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CreatedOrderReq
//
// Responses:
//  200: CreatedOrderResp

func CreatedOrdOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatedOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewCreatedOrdOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreatedOrdOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
