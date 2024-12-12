package project

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"btc_order/internal/logic/project"
	"btc_order/internal/svc"
	"btc_order/internal/types"
)

// swagger:route get /v1/project/list project GetProjectList
//
// 获取项目列表
//
// 获取项目列表
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ProjectListReq
//
// Responses:
//  200: ProjectListResp

func GetProjectListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProjectListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := project.NewGetProjectListLogic(r.Context(), svcCtx)
		resp, err := l.GetProjectList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
