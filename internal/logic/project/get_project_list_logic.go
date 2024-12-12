package project

import (
	"context"

	"btc_order/internal/svc"
	"btc_order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProjectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectListLogic {
	return &GetProjectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetProjectListLogic) GetProjectList(req *types.ProjectListReq) (resp *types.ProjectListResp, err error) {
	var projects []types.Project
	page := (req.Page - 1) * req.PageSize
	err = l.svcCtx.PgDB.ModelContext(l.ctx, &projects).
		Order("sort asc").
		Offset(int(page)).
		Limit(int(req.PageSize)).
		Select()
	if err != nil {
		return nil, err
	}
	count, err := l.svcCtx.PgDB.ModelContext(l.ctx, (*types.Project)(nil)).
		Count()
	if err != nil {
		return nil, err
	}
	resp = &types.ProjectListResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 200,
			Msg:  "success",
		}}

	for _, projectInfo := range projects {
		resp.Data.Data = append(resp.Data.Data, types.ProjectInfo{
			Name:        projectInfo.Name,
			Image:       projectInfo.Image,
			Type:        projectInfo.Type,
			Information: projectInfo.Information,
			Stage1:      projectInfo.Stage1,
			Stage2:      projectInfo.Stage2,
			Stage3:      projectInfo.Stage3,
		})
	}
	resp.Data.Total = uint64(count)
	return resp, nil
}
