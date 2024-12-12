package project

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/project"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/project/tick/{id} project GetProjectTickById
//

//

//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ProjectReq
//
// Responses:
//  200: ProjectTickeResp

func GetProjectTickByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProjectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := project.NewGetProjectTickByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetProjectTickById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
