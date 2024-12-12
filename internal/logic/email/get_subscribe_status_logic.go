package email

import (
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"context"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscribeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeStatusLogic {
	return &GetSubscribeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetSubscribeStatusLogic) GetSubscribeStatus(req *types.EmailReq) (resp *types.StatusResp, err error) {
	addr, err := url.QueryUnescape(req.Address)
	if err != nil {
		return nil, err
	}

	e := types.Email{}
	result, err := l.svcCtx.PgDB.ModelContext(l.ctx, &e).Where("address = ?", addr).Exists()
	if err != nil {
		return nil, err
	}

	return &types.StatusResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: result,
	}, nil
}
