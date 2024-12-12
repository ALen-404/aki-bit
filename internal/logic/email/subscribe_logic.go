package email

import (
	"context"
	"github.com/google/uuid"
	"time"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeLogic {
	return &SubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *SubscribeLogic) Subscribe(req *types.SubscribeReq) (resp *types.BaseMsgResp, err error) {
	e := types.Email{}

	result, err := l.svcCtx.PgDB.ModelContext(l.ctx, &e).Where("address = ?", req.Address).Exists()
	if err != nil {
		return nil, err
	}

	if result {
		return &types.BaseMsgResp{
			Code: 0,
			Msg:  "already subscribed",
		}, nil
	}

	e.ID = uuid.New()
	e.Address = req.Address
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err = l.svcCtx.PgDB.ModelContext(l.ctx, &e).Insert()
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{
		Code: 200,
		Msg:  "success",
	}, nil
}
