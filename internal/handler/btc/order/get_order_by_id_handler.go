package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/btc/order"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/btc/order/{id} btc_order GetOrderById
//
// 根据订单ID获取订单详情
//
// 根据订单ID获取订单详情
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: UUIDPathReq
//
// Responses:
//  200: OrderInfoResp

func GetOrderByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UUIDPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewGetOrderByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
