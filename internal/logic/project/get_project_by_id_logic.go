package project

import (
	"context"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProjectByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectByIdLogic {
	return &GetProjectByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetProjectByIdLogic) GetProjectById(req *types.ProjectReq) (resp *types.ProjecctResp, err error) {
	var projectInfo types.Project
	_, err = l.svcCtx.PgDB.ModelContext(l.ctx, &projectInfo).
		Where("id = ?", req.ID).Order("sort asc").
		Exists()
	if err != nil {
		return nil, err
	}
	resp = &types.ProjecctResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		},
		Data: types.ProjectInfo{
			Name:        projectInfo.Name,
			Image:       projectInfo.Image,
			Type:        projectInfo.Type,
			Information: projectInfo.Information,
			Stage1:      projectInfo.Stage1,
			Stage2:      projectInfo.Stage2,
			Stage3:      projectInfo.Stage3,
		},
	}
	return resp, nil
}
