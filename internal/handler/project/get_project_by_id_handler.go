package project

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/project"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/project/{id} project GetProjectById
//
// 根据tick获取项目详情
//
// 根据tick获取项目详情
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ProjectReq
//
// Responses:
//  200: ProjecctResp

func GetProjectByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProjectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := project.NewGetProjectByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetProjectById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
