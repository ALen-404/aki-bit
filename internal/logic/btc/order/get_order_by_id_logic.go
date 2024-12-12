package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/jsonx"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderByIdLogic {
	return &GetOrderByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetOrderByIdLogic) GetOrderById(req *types.UUIDPathReq) (resp *types.OrderInfoResp, err error) {
	btcOrder := types.BtcOrder{
		ID: uuid.MustParse(req.Id),
	}
	err = l.svcCtx.PgDB.ModelContext(l.ctx, &btcOrder).WherePK().Select()
	if err != nil {
		return nil, err
	}
	// 	files
	var files []types.File
	err = l.svcCtx.PgDB.ModelContext(l.ctx, &files).
		Where("btc_order_id = ?", req.Id).
		Order("index ASC").
		Select()

	var orderFiles []types.OrderFile
	for _, file := range files {
		orderFile := types.OrderFile{
			FileName: file.FileName,
			Size:     file.Size,
			TxId:     file.TxId,
			Address:  file.Address,
		}
		orderFiles = append(orderFiles, orderFile)
	}

	resp = &types.OrderInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
	}
	payAddress := ""
	if btcOrder.PayAddress != "" {
		var payAddresses []string
		err = jsonx.UnmarshalFromString(btcOrder.PayAddress, &payAddresses)
		if err != nil {
			return nil, err
		}
		payAddress = payAddresses[0]
	}
	boi := types.BtcOrderInfo{
		Id:             btcOrder.ID.String(),
		Type:           btcOrder.Type,
		TxId:           btcOrder.TxId,
		TxHash:         btcOrder.TxHash,
		Status:         btcOrder.Status,
		PayAddress:     payAddress,
		ReceiveAddress: btcOrder.ReceiveAddress,
		ClientId:       btcOrder.ClientId,
		Count:          btcOrder.Count,
		FeeRate:        btcOrder.FeeRate,
		MinerFee:       btcOrder.MinerFee,
		ServiceFee:     btcOrder.ServiceFee,
		Amount:         btcOrder.MinerFee + btcOrder.ServiceFee + btcOrder.Amount,
		ExTime:         btcOrder.ExTime,
		CreatedAt:      btcOrder.CreatedAt,
		Files:          orderFiles,
	}
	resp.Data = boi
	return resp, nil
}
