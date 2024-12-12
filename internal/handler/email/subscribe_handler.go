package email

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/email"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route post /v1/email/{address}/subscribe email Subscribe
//
// 邮件订阅
//
// 邮件订阅
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: EmailReq
//
// Responses:
//  200: BaseMsgResp

func SubscribeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubscribeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := email.NewSubscribeLogic(r.Context(), svcCtx)
		resp, err := l.Subscribe(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
