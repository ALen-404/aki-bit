package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/btc/order"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route put /v1/btc/order/{id} btc_order PutOrderById
//
// 根据订单ID修改订单信息
//
// 根据订单ID修改订单信息
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: OrderPutInfo
//
// Responses:
//  200: BaseMsgResp

func PutOrderByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderPutInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewPutOrderByIdLogic(r.Context(), svcCtx)
		resp, err := l.PutOrderById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
