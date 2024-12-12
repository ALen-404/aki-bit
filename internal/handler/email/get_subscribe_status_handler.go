package email

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/email"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/email/{address}/status email GetSubscribeStatus
//
// 获取邮件订阅状态
//
// 获取邮件订阅状态
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: EmailReq
//
// Responses:
//  200: StatusResp

func GetSubscribeStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := email.NewGetSubscribeStatusLogic(r.Context(), svcCtx)
		resp, err := l.GetSubscribeStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
