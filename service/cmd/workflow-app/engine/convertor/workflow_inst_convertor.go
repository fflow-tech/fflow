package convertor

import (
	"fmt"
	"sort"
	"time"

	pb "github.com/fflow-tech/fflow/api/workflow-app/engine"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	// InstConvertor 转换器
	InstConvertor = &instConvertorImpl{}
)

type instConvertorImpl struct {
}

// ConvertDTOToPb 转换
func (*instConvertorImpl) ConvertDTOToPb(oldDetail *dto.WorkflowInstDTO, newDetail *pb.GetInstDetailRsp) error {
	newDetail.DefID = oldDetail.WorkflowDef.DefID
	newDetail.InstID = oldDetail.InstID
	newDetail.Name = oldDetail.Name
	newDetail.Status = oldDetail.Status.String()
	newDetail.Input = utils.MapToStr(oldDetail.Input)
	newDetail.Creator = oldDetail.Creator
	newDetail.StartAt = oldDetail.StartAt.Unix()
	newDetail.CompletedAt = oldDetail.CompletedAt.Unix()
	newDetail.CostTime = utils.GetCostTime(newDetail.StartAt, newDetail.CompletedAt)
	newDetail.CurNodes = oldDetail.ExecutePath[len(oldDetail.ExecutePath)-1]

	nodeItems, err := getNodeItems(oldDetail)
	if err != nil {
		return err
	}
	newDetail.NodeItems = nodeItems

	return nil
}

// getNodeItems 获取所有的节点实例
func getNodeItems(oldDetail *dto.WorkflowInstDTO) ([]*pb.NodeItem, error) {
	schedNodeItems, err := getSchedNodeItems(oldDetail)
	if err != nil {
		return nil, err
	}

	notSchedNodeItems, err := getNotSchedNodeItems(oldDetail)
	if err != nil {
		return nil, err
	}

	// 合并所有的 nodeItem 并按照 seq 排序
	nodeItems := append(schedNodeItems, notSchedNodeItems...)
	sort.Slice(nodeItems, func(i, j int) bool {
		return nodeItems[i].Seq < nodeItems[j].Seq
	})

	return nodeItems, nil
}

// getSchedNodeItems 获取调度了的节点实例
func getSchedNodeItems(oldDetail *dto.WorkflowInstDTO) ([]*pb.NodeItem, error) {
	schedNodeItems := []*pb.NodeItem{}
	for _, nodeInst := range oldDetail.SchedNodeInsts {
		nodeItem, err := NodeInstConvertor.ConvertDTOToPb(&dto.NodeInstDTO{NodeInst: *nodeInst})
		if err != nil {
			return nil, err
		}
		schedNodeItems = append(schedNodeItems, nodeItem)
	}
	return schedNodeItems, nil
}

// getNotSchedNodeItems 获取未调度的节点实例默认值
func getNotSchedNodeItems(oldDetail *dto.WorkflowInstDTO) ([]*pb.NodeItem, error) {
	notSchedNodeItems := []*pb.NodeItem{}
	for _, node := range oldDetail.WorkflowDef.Nodes {
		for nodeRefName := range node {
			// 如果已经调度过的则跳过
			if isInSchedNodeItems(nodeRefName, oldDetail.SchedNodeInsts) {
				continue
			}
			// 否则创建一个空的节点实例并添加到未调度的节点中
			newNodeInst, err := createNotSchedNodeInst(&oldDetail.WorkflowInst, nodeRefName)
			if err != nil {
				return nil, err
			}
			newNodeItem, err := NodeInstConvertor.ConvertDTOToPb(&dto.NodeInstDTO{NodeInst: *newNodeInst})
			if err != nil {
				return nil, err
			}
			notSchedNodeItems = append(notSchedNodeItems, newNodeItem)
		}
	}
	return notSchedNodeItems, nil
}

// createNotSchedNodeInst 获取未调度节点实例
func createNotSchedNodeInst(inst *entity.WorkflowInst, nodeRefName string) (*entity.NodeInst, error) {
	nodeInst, err := entity.NewNodeInstByNodeRefName(*inst, nodeRefName)
	if err != nil {
		return nil, err
	}
	// 未调度的节点状态为空
	nodeInst.Status = entity.NodeInstStatus{}
	nodeInst.ScheduledAt = time.Time{}

	return nodeInst, nil
}

// isInSchedNodeItems 根据 refName 查找节点实例是否已经在调度过的 nodeItems 中了
func isInSchedNodeItems(refName string, nodeItems []*entity.NodeInst) bool {
	for _, nodeItem := range nodeItems {
		if refName == nodeItem.BasicNodeDef.RefName {
			return true
		}
	}
	return false
}

// buildURL 根据请求的 input 和 output 获取蓝盾构建地址的 URL
// 只是为了兼容, 所以不处理错误
func buildURL(input string, output string) string {
	if input == "" || output == "" {
		return ""
	}
	inputBytes := []byte(input)
	outputBytes := []byte(output)

	projectID, err := utils.GetStrFromJson(inputBytes, "projectID")
	if err != nil {
		return ""
	}
	pipelineID, err := utils.GetStrFromJson(inputBytes, "pipelineID")
	if err != nil {
		return ""
	}
	buildID, err := utils.GetStrFromJson(outputBytes, "data", "id")
	if err != nil {
		return ""
	}
	// 任何一个字段为空时则无法得到蓝盾的地址直接返回
	if projectID == "" || pipelineID == "" || buildID == "" {
		return ""
	}

	return fmt.Sprintf("https://devops.woa.com/console/pipeline/%s/%s/detail/%s", projectID, pipelineID, buildID)
}

// ConvertDTOListToPb 转换
func (c *instConvertorImpl) ConvertDTOListToPb(oldDetails []*dto.WorkflowInstDTO, total int64,
	newDetails *pb.GetInstListRsp) error {
	newDetails.InstDetails = []*pb.InstDetail{}
	for _, oldDetail := range oldDetails {
		newDetail := &pb.InstDetail{
			Name:        oldDetail.Name,
			InstID:      oldDetail.InstID,
			Status:      oldDetail.Status.String(),
			Creator:     oldDetail.Creator,
			CreateAt:    oldDetail.StartAt.Unix(),
			StartAt:     oldDetail.StartAt.Unix(),
			CompletedAt: oldDetail.CompletedAt.Unix(),
		}
		newDetails.InstDetails = append(newDetails.InstDetails, newDetail)
	}
	newDetails.Total = int32(total)
	return nil
}

// ConvertStartPbToDTO 转换
func (*instConvertorImpl) ConvertStartPbToDTO(req *pb.StartInstReq) (*dto.StartWorkflowInstDTO, error) {
	createInstDTO := &dto.StartWorkflowInstDTO{}
	createInstDTO.DefID = req.DefID
	createInstDTO.Namespace = req.BasicReq.Namespace
	createInstDTO.Creator = req.BasicReq.Operator
	createInstDTO.Reason = req.BasicReq.Reason
	createInstDTO.Name = req.InstName
	createInstDTO.DebugMode = req.DebugMode
	var err error
	if createInstDTO.Input, err = utils.JsonStrToMap(req.Input); err != nil {
		return nil, err
	}

	return createInstDTO, nil
}

// ConvertGetPbToDTO 转换
func (*instConvertorImpl) ConvertGetPbToDTO(req *pb.GetInstDetailReq) (*dto.GetWorkflowInstDTO, error) {
	return &dto.GetWorkflowInstDTO{
		InstID: req.InstID,
	}, nil
}

// ConvertGetListPbToDTO 转换
func (*instConvertorImpl) ConvertGetListPbToDTO(req *pb.GetInstListReq) (*dto.GetWorkflowInstListDTO, error) {
	return &dto.GetWorkflowInstListDTO{
		Namespace: req.BasicReq.Namespace,
		PageQuery: constants.NewPageQuery(int(req.PageIndex), int(req.PageSize)),
		Status:    entity.GetInstStatusByStrValue(req.Status),
		DefID:     req.DefID,
		AscOrder:  req.AscOrder,
	}, nil
}

// ConvertWebGetListToDTO 转换
func (*instConvertorImpl) ConvertWebGetListToDTO(req *dto.GetWebWorkflowInstListDTO) *dto.GetWorkflowInstListDTO {
	if req.PageQuery == nil {
		req.PageQuery = constants.NewDefaultPageQuery()
	}

	return &dto.GetWorkflowInstListDTO{
		Namespace: req.Namespace,
		PageQuery: constants.NewPageQuery(req.PageIndex, req.PageSize),
		Operator:  req.Operator,
		Name:      req.Name,
		Status:    entity.GetInstStatusByStrValue(req.Status),
		DefID:     req.DefID,
		AscOrder:  req.AscOrder,
	}
}

// ConvertRestartPbToDTO 转换
func (*instConvertorImpl) ConvertRestartPbToDTO(req *pb.RestartInstReq) (*dto.RestartWorkflowInstDTO, error) {
	restartInstDTO := &dto.RestartWorkflowInstDTO{}
	restartInstDTO.InstID = req.InstID
	restartInstDTO.NodeRefName = req.NodeRefName
	restartInstDTO.Namespace = req.BasicReq.Namespace
	restartInstDTO.Operator = req.BasicReq.Operator
	restartInstDTO.Reason = req.BasicReq.Reason
	var err error
	if restartInstDTO.Input, err = utils.JsonStrToMap(req.Input); err != nil {
		return nil, err
	}

	return restartInstDTO, nil
}

// ConvertCancelPbToDTO 转换
func (*instConvertorImpl) ConvertCancelPbToDTO(req *pb.CancelInstReq) (*dto.CancelWorkflowInstDTO, error) {
	cancelInstDTO := &dto.CancelWorkflowInstDTO{}
	cancelInstDTO.Namespace = req.BasicReq.Namespace
	cancelInstDTO.Operator = req.BasicReq.Operator
	cancelInstDTO.Reason = req.BasicReq.Reason
	cancelInstDTO.InstID = req.InstID
	return cancelInstDTO, nil
}

// ConvertCompletePbToDTO 转换
func (*instConvertorImpl) ConvertCompletePbToDTO(req *pb.CompleteInstReq) (*dto.CompleteWorkflowInstDTO, error) {
	completeDTO := &dto.CompleteWorkflowInstDTO{}
	completeDTO.Namespace = req.BasicReq.Namespace
	completeDTO.Operator = req.BasicReq.Operator
	completeDTO.Reason = req.BasicReq.Reason
	completeDTO.Status = entity.GetInstStatusByStrValue(req.Status)
	completeDTO.InstID = req.InstID
	return completeDTO, nil
}

// ConvertUpdateCtxPbToDTO 转换
func (*instConvertorImpl) ConvertUpdateCtxPbToDTO(req *pb.UpdateInstContextReq) (*dto.UpdateWorkflowInstCtxDTO, error) {
	updateCtxDTO := &dto.UpdateWorkflowInstCtxDTO{}
	updateCtxDTO.Namespace = req.BasicReq.Namespace
	updateCtxDTO.Operator = req.BasicReq.Operator
	updateCtxDTO.Reason = req.BasicReq.Reason
	var err error
	if updateCtxDTO.Context, err = utils.JsonStrToMap(req.Context); err != nil {
		return nil, err
	}
	updateCtxDTO.InstID = req.InstID
	return updateCtxDTO, nil
}

// ConvertPausePbToDTO 转换
func (*instConvertorImpl) ConvertPausePbToDTO(req *pb.PauseInstReq) (*dto.PauseWorkflowInstDTO, error) {
	pauseInstDTO := &dto.PauseWorkflowInstDTO{}
	pauseInstDTO.Namespace = req.BasicReq.Namespace
	pauseInstDTO.Operator = req.BasicReq.Operator
	pauseInstDTO.Reason = req.BasicReq.Reason
	pauseInstDTO.InstID = req.InstID
	return pauseInstDTO, nil
}

// ConvertResumePbToDTO 转换
func (*instConvertorImpl) ConvertResumePbToDTO(req *pb.ResumeInstReq) (*dto.ResumeWorkflowInstDTO, error) {
	resumeInstDTO := &dto.ResumeWorkflowInstDTO{}
	resumeInstDTO.InstID = req.InstID
	resumeInstDTO.Namespace = req.BasicReq.Namespace
	resumeInstDTO.Operator = req.BasicReq.Operator
	resumeInstDTO.Reason = req.BasicReq.Reason
	return resumeInstDTO, nil
}
