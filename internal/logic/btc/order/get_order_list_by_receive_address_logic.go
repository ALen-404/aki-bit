package order

import (
	"context"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListByReceiveAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListByReceiveAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListByReceiveAddressLogic {
	return &GetOrderListByReceiveAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetOrderListByReceiveAddressLogic) GetOrderListByReceiveAddress(req *types.OrderListReq) (resp *types.OrderListResp, err error) {
	var btcOrders []types.BtcOrder

	page := (req.Page - 1) * req.PageSize

	err = l.svcCtx.PgDB.ModelContext(l.ctx, &btcOrders).Column("id", "created_at", "status").
		Where("receive_address = ?", req.ReceiveAddress).
		Order("created_at desc").
		Offset(int(page)).
		Limit(int(req.PageSize)).
		Select()
	if err != nil {
		return nil, err
	}
	count, err := l.svcCtx.PgDB.ModelContext(l.ctx, (*types.BtcOrder)(nil)).
		Where("receive_address = ?", req.ReceiveAddress).
		Count()
	if err != nil {
		return nil, err
	}
	resp = &types.OrderListResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		}}
	for _, bo := range btcOrders {
		resp.Data.Data = append(resp.Data.Data, types.BtcOrderListInfo{
			ID:        bo.ID.String(),
			CreatedAt: bo.CreatedAt,
			Status:    bo.Status,
		})
	}
	resp.Data.Total = uint64(count)
	return resp, nil
}
