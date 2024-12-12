package project

import (
	"context"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProjectTickByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectTickByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectTickByIdLogic {
	return &GetProjectTickByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetProjectTickByIdLogic) GetProjectTickById(req *types.ProjectReq) (resp *types.ProjectTickeResp, err error) {
	var projectInfo types.Project
	_, err = l.svcCtx.PgDB.ModelContext(l.ctx, &projectInfo).
		Where("id = ?", req.ID).
		Exists()
	if err != nil {
		return nil, err
	}
	resp = &types.ProjectTickeResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: types.ProjectTick{
			Tick: projectInfo.Tick,
		},
	}
	return resp, nil
}
