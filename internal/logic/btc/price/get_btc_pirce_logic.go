package price

import (
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBtcPirceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBtcPirceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBtcPirceLogic {
	return &GetBtcPirceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetBtcPirceLogic) GetBtcPirce() (resp *types.BtcPriceResp, err error) {
	resp = &types.BtcPriceResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
	}

	priceResp, err := http.Get("https://www.geniidata.com/api/uda/flow/btc/price")
	if err != nil {
		return nil, err
	}
	defer priceResp.Body.Close()
	var body map[string]interface{}
	err = json.NewDecoder(priceResp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	content := body["data"]
	resp.Data = content
	return resp, nil
}
