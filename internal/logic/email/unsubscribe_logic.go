package email

import (
	"btc_order/internal/svc"
	"btc_order/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsubscribeLogic {
	return &UnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *UnsubscribeLogic) Unsubscribe(req *types.SubscribeReq) (resp *types.BaseMsgResp, err error) {
	e := types.Email{}

	result, err := l.svcCtx.PgDB.ModelContext(l.ctx, &e).Where("address = ?", req.Address).Exists()
	if err != nil {
		return nil, err
	}

	if !result {
		return &types.BaseMsgResp{
			Code: 0,
			Msg:  "already unsubscribed",
		}, nil
	}

	_, err = l.svcCtx.PgDB.ModelContext(l.ctx, &e).Where("address = ?", req.Address).Delete()
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{
		Code: 200,
		Msg:  "success",
	}, nil
}
