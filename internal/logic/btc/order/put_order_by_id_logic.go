package order

import (
	"context"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutOrderByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutOrderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutOrderByIdLogic {
	return &PutOrderByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *PutOrderByIdLogic) PutOrderById(req *types.OrderPutInfoReq) (resp *types.BaseMsgResp, err error) {
	tx, err := l.svcCtx.PgDB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Close()
	}()
	// 	files
	var files []types.File
	err = tx.ModelContext(l.ctx, &files).
		Where("btc_order_id = ?", req.Id).
		Order("index ASC").
		Select()
	if err != nil {
		return nil, err
	}
	for i, file := range files {
		file.TxId = req.RevealTxs[i]
		_, err = tx.ModelContext(l.ctx, &file).WherePK().Update()
		if err != nil {
			return nil, err
		}
	}

	orderInfo := map[string]interface{}{
		"tx_id":       req.CommitTx,
		"pay_address": req.CommitAddrs,
	}
	_, err = tx.ModelContext(l.ctx, &orderInfo).TableExpr("btc_order").Where("id = ?", req.Id).Update()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Code: 200, Msg: "success"}, nil
}
