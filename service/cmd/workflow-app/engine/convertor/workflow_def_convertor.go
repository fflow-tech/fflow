package convertor

import (
	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/jinzhu/copier"
)

var (
	// DefConvertor 转换器
	DefConvertor = &defConvertorImpl{}
)

type defConvertorImpl struct {
}

// ConvertCreatePbToDTO 转换
func (*defConvertorImpl) ConvertCreatePbToDTO(req *pb.CreateDefReq) (*dto.CreateWorkflowDefDTO, error) {
	createDefDTO := &dto.CreateWorkflowDefDTO{}
	err := copier.Copy(createDefDTO, req)
	if err != nil {
		log.Errorf("Failed to copy error, caused by %s, req:%s", err, utils.StructToJsonStr(req))
		return nil, err
	}

	createDefDTO.Namespace = req.BasicReq.Namespace
	createDefDTO.DefJson = req.Content
	createDefDTO.Creator = req.Author
	createDefDTO.Description = req.Desc
	return createDefDTO, nil
}

// ConvertUpdatePbToDTO 转换
func (*defConvertorImpl) ConvertUpdatePbToDTO(req *pb.UpdateDefReq) (*dto.CreateWorkflowDefDTO, error) {
	updateDefDTO := &dto.CreateWorkflowDefDTO{}
	err := copier.Copy(updateDefDTO, req)
	if err != nil {
		log.Errorf("Failed to copy error, caused by %s, req:%v", err, req)
		return nil, err
	}
	updateDefDTO.DefID = req.DefID
	updateDefDTO.Namespace = req.BasicReq.Namespace
	updateDefDTO.DefJson = req.Content
	updateDefDTO.Creator = req.Author
	updateDefDTO.Description = req.Desc
	return updateDefDTO, nil
}

// ConvertGetPbToDTO 转换
func (*defConvertorImpl) ConvertGetPbToDTO(req *pb.GetDefDetailReq) (*dto.GetWorkflowDefDTO, error) {
	return &dto.GetWorkflowDefDTO{
		Namespace: req.BasicReq.Namespace,
		DefID:     req.DefID,
	}, nil
}

// ConvertEnablePbToDTO 转换
func (*defConvertorImpl) ConvertEnablePbToDTO(req *pb.EnableDefReq) (*dto.EnableWorkflowDefDTO, error) {
	return &dto.EnableWorkflowDefDTO{
		Namespace: req.BasicReq.Namespace,
		DefID:     req.DefID,
	}, nil
}

// ConvertDisablePbToDTO 转换
func (*defConvertorImpl) ConvertDisablePbToDTO(req *pb.DisableDefReq) (*dto.DisableWorkflowDefDTO, error) {
	return &dto.DisableWorkflowDefDTO{
		Namespace: req.BasicReq.Namespace,
		DefID:     req.DefID,
	}, nil
}

// ConvertDTOToPb 转换
func (*defConvertorImpl) ConvertDTOToPb(d *dto.WorkflowDefDTO) *pb.WorkflowDef {
	return &pb.WorkflowDef{
		Name:     d.Name,
		Desc:     d.Description,
		Author:   d.Creator,
		Format:   "json",
		Content:  d.DefJson,
		CreateAt: d.CreatedAt.Unix(),
		UpdateAt: d.CreatedAt.Unix(),
		Status:   d.Status.String(),
	}
}

// ConvertWebGetListToDTO 转换
func (*defConvertorImpl) ConvertWebGetListToDTO(req *dto.PageQueryWebWorkflowDefDTO) *dto.PageQueryWorkflowDefDTO {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	return &dto.PageQueryWorkflowDefDTO{
		Namespace: req.Namespace,
		Operator:  req.Operator,
		DefID:     req.DefID,
		Name:      req.Name,
		Status:    entity.GetDefStatusByStrValue(req.Status).IntValue(),
		Version:   req.Version,
		GroupBy:   req.GroupBy,
		PageQuery: req.PageQuery,
		Order:     req.Order,
	}
}
