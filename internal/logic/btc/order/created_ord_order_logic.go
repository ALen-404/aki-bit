package order

import (
	"context"
	"github.com/google/uuid"
	"math"
	"strings"
	"time"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatedOrdOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatedOrdOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedOrdOrderLogic {
	return &CreatedOrdOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *CreatedOrdOrderLogic) CreatedOrdOrder(req *types.CreatedOrderReq) (resp *types.CreatedOrderResp, err error) {
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
	now := time.Now()
	exTime := now.Add(30 * time.Minute)
	orderId := uuid.New()

	fileSize := 0
	for i, file := range req.Files {
		s := len(file.FileName)
		strings.Split(file.DataUrl, "")
		fileSize += s
		_, err = l.svcCtx.PgDB.ModelContext(l.ctx, &types.File{
			ID:         uuid.New(),
			BtcOrderId: orderId,
			FileName:   file.FileName,
			DataUrl:    file.DataUrl,
			Size:       uint64(s),
			Index:      uint64(i),
		}).Insert()
		if err != nil {
			return nil, err
		}
	}

	af, nf, sf := calculateFee(req.ReceiveAddress, float64(req.Amount),
		float64(len(req.Files)), float64(fileSize), float64(len("text/plain;charset=utf-8")), float64(req.FeeRate))
	btcOrder := &types.BtcOrder{
		ID:             orderId,
		Type:           1,
		Status:         "pending",
		ReceiveAddress: req.ReceiveAddress,
		ClientId:       req.ClientId,
		FeeRate:        req.FeeRate,
		Amount:         af,
		CreatedAt:      now.UnixMilli(),
		ExTime:         exTime.UnixMilli(),
		Count:          uint64(len(req.Files)),
		MinerFee:       nf,
		ServiceFee:     sf,
	}

	_, err = l.svcCtx.PgDB.ModelContext(l.ctx, btcOrder).Insert()
	if err != nil {
		return nil, err
	}

	bo := &types.BtcOrder{
		ID: orderId,
	}
	err = l.svcCtx.PgDB.ModelContext(l.ctx, bo).WherePK().Select()
	if err != nil {
		return nil, err
	}

	resp = &types.CreatedOrderResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: types.CreatedOrderInfo{
			Id:             bo.ID.String(),
			ClientId:       bo.ClientId,
			Type:           bo.Type,
			Status:         bo.Status,
			ReceiveAddress: bo.ReceiveAddress,
			Amount:         nf + sf + af,
			Count:          bo.Count,
			FeeRate:        bo.FeeRate,
			MinerFee:       bo.MinerFee,
			ServiceFee:     bo.ServiceFee,
			ExTime:         bo.ExTime,
			CreatedAt:      bo.CreatedAt,
		},
	}
	return resp, nil
}

func calculateFee(address string, inscriptionBalance, fileCount, fileSize, contentTypeSize, feeRate float64) (amount, networkFee, serviceFee uint64) {
	var addrSize float64 = 25 + 1 // p2pkh
	var baseSize float64 = 88
	if strings.Contains(address, "bc1q") || strings.Contains(address, "tb1q") {
		addrSize = 22 + 1
	} else if strings.Contains(address, "bc1p") || strings.Contains(address, "tb1p") {
		addrSize = 34 + 1
	} else if strings.Contains(address, "2") || strings.Contains(address, "3") {
		addrSize = 23 + 1
	}
	balance := inscriptionBalance * fileCount

	networkSats := math.Ceil(((fileSize+contentTypeSize)/4 + (baseSize + 8 + addrSize + 8 + 23)) * feeRate)
	if fileCount > 1 {
		networkSats = math.Ceil(((fileSize+contentTypeSize)/4 + (baseSize + 8 + addrSize + (35+8)*(fileCount-1) + 8 + 23 + (baseSize+8+addrSize+0.5)*(fileCount-1))) * feeRate)
	}

	return uint64(balance), uint64(networkSats), 0
}
