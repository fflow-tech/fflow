package convertor

import (
	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	// NodeInstConvertor 转换器
	NodeInstConvertor = &nodeInstConvertor{}
)

type nodeInstConvertor struct {
}

// ConvertGetPbToDTO 转换
func (*nodeInstConvertor) ConvertGetPbToDTO(req *pb.GetNodeInstDetailReq) (*dto.GetNodeInstDTO, error) {
	return &dto.GetNodeInstDTO{
		NodeInstID: req.NodeInstID,
	}, nil
}

// ConvertDTOToPb 转换
func (*nodeInstConvertor) ConvertDTOToPb(nodeInstDTO *dto.NodeInstDTO) (*pb.NodeItem, error) {
	nodeItem := &pb.NodeItem{
		InstID:      nodeInstDTO.InstID,
		Status:      nodeInstDTO.Status.String(),
		Input:       utils.MapToStr(nodeInstDTO.Input),
		Output:      utils.MapToStr(nodeInstDTO.Output),
		PollInput:   utils.MapToStr(nodeInstDTO.PollInput),
		PollOutput:  utils.MapToStr(nodeInstDTO.PollOutput),
		NodeRefName: nodeInstDTO.BasicNodeDef.RefName,
		Name:        nodeInstDTO.BasicNodeDef.Name,
		Type:        nodeInstDTO.BasicNodeDef.Type.String(),
		Biz:         utils.MapToStr(nodeInstDTO.Biz),
		StartAt:     nodeInstDTO.ScheduledAt.Unix(),
		CompletedAt: nodeInstDTO.CompletedAt.Unix(),
		CostTime:    utils.GetCostTime(nodeInstDTO.ScheduledAt.Unix(), nodeInstDTO.CompletedAt.Unix()),
		Seq:         int32(nodeInstDTO.BasicNodeDef.Index),
		Nexts:       nodeInstDTO.Nexts,
		Parents:     nodeInstDTO.Parents,
	}
	// 为了兼容老版本的引擎，对于蓝盾相关节点根据接口入参和输出获取蓝盾的构建地址
	nodeItem.Url = buildURL(nodeItem.Input, nodeItem.Output)
	return nodeItem, nil
}

// ConvertRerunPbToDTO 转换
func (*nodeInstConvertor) ConvertRerunPbToDTO(req *pb.RerunNodeReq) (*dto.RerunNodeDTO, error) {
	rerunDTO := &dto.RerunNodeDTO{}
	rerunDTO.InstID = req.InstID
	rerunDTO.Namespace = req.BasicReq.Namespace
	rerunDTO.NodeRefName = req.NodeRefName
	rerunDTO.Reason = req.BasicReq.Reason
	return rerunDTO, nil
}

// ConvertResumePbToDTO 转换
func (*nodeInstConvertor) ConvertResumePbToDTO(req *pb.ResumeNodeReq) (*dto.ResumeNodeDTO, error) {
	resumeNode := &dto.ResumeNodeDTO{}
	resumeNode.InstID = req.InstID
	resumeNode.NodeInstID = req.NodeInstID
	resumeNode.Namespace = req.BasicReq.Namespace
	resumeNode.NodeRefName = req.NodeRefName
	resumeNode.Reason = req.BasicReq.Reason
	return resumeNode, nil
}

// ConvertCancelPbToDTO 转换
func (*nodeInstConvertor) ConvertCancelPbToDTO(req *pb.CancelNodeReq) (*dto.CancelNodeDTO, error) {
	cancelDTO := &dto.CancelNodeDTO{}
	cancelDTO.InstID = req.InstID
	cancelDTO.NodeInstID = req.NodeInstID
	cancelDTO.Namespace = req.BasicReq.Namespace
	cancelDTO.NodeRefName = req.NodeRefName
	cancelDTO.Reason = req.BasicReq.Reason
	return cancelDTO, nil
}

// ConvertCompletePbToDTO 转换
func (*nodeInstConvertor) ConvertCompletePbToDTO(req *pb.CompleteNodeReq) (*dto.CompleteNodeDTO, error) {
	completeDTO := &dto.CompleteNodeDTO{}
	completeDTO.InstID = req.InstID
	completeDTO.NodeInstID = req.NodeInstID
	completeDTO.Status = entity.GetNodeInstStatusByStrValue(req.Status)
	completeDTO.Namespace = req.BasicReq.Namespace
	completeDTO.NodeRefName = req.NodeRefName
	completeDTO.Reason = req.BasicReq.Reason
	if req.FailedReason != "" {
		completeDTO.Reason = req.FailedReason
	}
	var err error
	if completeDTO.Output, err = utils.JsonStrToMap(req.NodeOutput); err != nil {
		return nil, err
	}
	return completeDTO, nil
}

// ConvertSkipPbToDTO 转换
func (*nodeInstConvertor) ConvertSkipPbToDTO(req *pb.SkipNodeReq) (*dto.SkipNodeDTO, error) {
	skipDTO := &dto.SkipNodeDTO{}
	skipDTO.InstID = req.InstID
	skipDTO.Namespace = req.BasicReq.Namespace
	skipDTO.NodeRefName = req.NodeRefName
	skipDTO.Reason = req.BasicReq.Reason
	return skipDTO, nil
}

// ConvertCancelSkipPbToDTO 转换
func (*nodeInstConvertor) ConvertCancelSkipPbToDTO(req *pb.CancelSkipNodeReq) (*dto.CancelSkipNodeDTO, error) {
	cancelSkipDTO := &dto.CancelSkipNodeDTO{}
	cancelSkipDTO.InstID = req.InstID
	cancelSkipDTO.Namespace = req.BasicReq.Namespace
	cancelSkipDTO.NodeRefName = req.NodeRefName
	cancelSkipDTO.Reason = req.BasicReq.Reason
	return cancelSkipDTO, nil
}
