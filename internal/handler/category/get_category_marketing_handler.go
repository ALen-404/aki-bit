package category

import (
	"btc_order/internal/logic/category"
	"btc_order/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route get /v1/market/category/latest category GetCategoryMarketing
//
// 获取BRC20市场总体信息
//
// 获取BRC20市场总体信息
//
// Responses:
//  200: CategoryMarketInfoResp

func GetCategoryMarketingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := category.NewGetCategoryMarketingLogic(r.Context(), svcCtx)
		resp, err := l.GetCategoryMarketing()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
