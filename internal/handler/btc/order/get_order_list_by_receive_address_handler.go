package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/btc/order"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/btc/order/list/{receiveAddress} btc_order GetOrderListByReceiveAddress
//
// 根据钱包地址获取订单列表
//
// 根据钱包地址获取订单列表
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: OrderListReq
//
// Responses:
//  200: OrderListResp

func GetOrderListByReceiveAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewGetOrderListByReceiveAddressLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderListByReceiveAddress(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
