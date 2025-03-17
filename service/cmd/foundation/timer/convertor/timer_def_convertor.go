// Package convertor 负责 PB->DTO 转换
package convertor

import (
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/jinzhu/copier"

	pb "github.com/fflow-tech/fflow/api/foundation/timer"
)

var (
	// TimerDefConvertor 定时器转换器
	TimerDefConvertor = &timerDefConvertor{}
)

type timerDefConvertor struct {
}

// ConvertCreatePbToDTO 创建PB转换DTO
func (*timerDefConvertor) ConvertCreatePbToDTO(req *pb.CreateTimerReq) (*dto.CreateTimerDefDTO, error) {
	createDefDTO := &dto.CreateTimerDefDTO{}
	err := copier.Copy(createDefDTO, req)
	if err != nil {
		log.Errorf("Failed to ConvertCreatePbToDTO copy error, caused by %s, req:%s", err,
			utils.StructToJsonStr(req))
		return nil, err
	}
	return createDefDTO, nil
}

// ConvertEnablePbToDTO 激活PB转换DTO
func (*timerDefConvertor) ConvertEnablePbToDTO(req *pb.EnableTimerReq) (*dto.ChangeTimerStatusDTO, error) {
	changeDTO := &dto.ChangeTimerStatusDTO{}
	changeDTO.DefID = req.DefId
	changeDTO.Status = entity.Enabled.ToInt()
	return changeDTO, nil
}

// ConvertDisablePbToDTO 去激活PB转换DTO
func (*timerDefConvertor) ConvertDisablePbToDTO(req *pb.DisableTimerReq) (*dto.ChangeTimerStatusDTO, error) {
	changeDTO := &dto.ChangeTimerStatusDTO{}
	changeDTO.DefID = req.DefId
	changeDTO.Status = entity.Disabled.ToInt()
	return changeDTO, nil
}

// ConvertGetDefPbToDTO 获取详情PB转换DTO
func (*timerDefConvertor) ConvertGetDefPbToDTO(req *pb.GetTimerDefReq) (*dto.GetTimerDefDTO, error) {
	getDTO := &dto.GetTimerDefDTO{}
	getDTO.DefID = req.DefId
	return getDTO, nil
}

// ConvertGetDefListPbToDTO 获取定义列表PB转换DTO
func (*timerDefConvertor) ConvertGetDefListPbToDTO(req *pb.GetTimerDefListReq) (*dto.PageQueryTimeDefDTO, error) {
	getDTO := &dto.PageQueryTimeDefDTO{}
	err := copier.Copy(getDTO, req)
	if err != nil {
		log.Errorf("Failed to ConvertGetDefPbToDTO copy error, caused by %s, req:%s", err,
			utils.StructToJsonStr(req))
		return nil, err
	}
	return getDTO, nil
}

// ConvertDelDefPbToDTO 删除定义PB转换DTO
func (*timerDefConvertor) ConvertDelDefPbToDTO(req *pb.DeleteTimerDefReq) (*dto.DeleteTimerDefDTO, error) {
	return &dto.DeleteTimerDefDTO{
		DefID: req.DefId,
		App:   req.App,
		Name:  req.Name,
	}, nil
}

// ConvertGetRunHistoryListPbToDTO 获取执行历史列表PB转换DTO
func (*timerDefConvertor) ConvertGetRunHistoryListPbToDTO(req *pb.GetRunHistoryListReq) (*dto.PageQueryRunHistoryDTO,
	error) {
	getHistoryListDTO := &dto.PageQueryRunHistoryDTO{}
	getHistoryListDTO.DefID = req.DefId
	getHistoryListDTO.PageQuery = constants.NewPageQuery(int(req.PageIndex), int(req.PageSize))
	getHistoryListDTO.Order = constants.NewDefaultOrder()
	if req.AscOrder {
		getHistoryListDTO.Order.Order = constants.DescOrder
	}

	return getHistoryListDTO, nil
}

// ConvertGetRunHistoryRspListDTOToPB 获取执行历史列表返回DTO转换PB
func (t *timerDefConvertor) ConvertGetRunHistoryRspListDTOToPB(d []*dto.GetRunHistoryRspDTO) ([]*pb.TaskRunHistory,
	error) {
	taskRunHistories := make([]*pb.TaskRunHistory, 0, len(d))
	for _, runHistoryDTO := range d {
		taskRunHistory, err := t.ConvertGetRunHistoryRspDTOToPB(runHistoryDTO)
		if err != nil {
			return nil, err
		}
		taskRunHistories = append(taskRunHistories, taskRunHistory)
	}
	return taskRunHistories, nil
}

// ConvertGetRunHistoryRspDTOToPB 获取执行历史返回DTO转换PB
func (*timerDefConvertor) ConvertGetRunHistoryRspDTOToPB(d *dto.GetRunHistoryRspDTO) (*pb.TaskRunHistory,
	error) {
	rsp := &pb.TaskRunHistory{}
	err := copier.Copy(rsp, d)
	if err != nil {
		log.Errorf("Failed to copy error, caused by %s, req:%s", err, utils.StructToJsonStr(d))
		return nil, err
	}
	rsp.DefId = d.DefID
	return rsp, nil
}

// ConvertDefDTOToTimerDefPB 获取定时器定义DTO转换PB
func (*timerDefConvertor) ConvertDefDTOToTimerDefPB(defDTO *dto.TimerDefDTO) (*pb.TimerDef, error) {
	timerDef := &pb.TimerDef{}
	err := copier.Copy(timerDef, defDTO)
	if err != nil {
		log.Errorf("Failed to copy error, caused by %s, req:%s", err, utils.StructToJsonStr(defDTO))
		return nil, err
	}
	timerDef.DefId = defDTO.DefID
	timerDef.Status = uint32(defDTO.Status)
	return timerDef, nil
}

// ConvertDefListDTOToTimerDefPB 获取定时器定义列表DTO转换PB
func (t *timerDefConvertor) ConvertDefListDTOToTimerDefPB(defDTOs []*dto.TimerDefDTO) ([]*pb.TimerDef, error) {
	timerDefList := make([]*pb.TimerDef, 0, len(defDTOs))
	for _, defDTO := range defDTOs {
		timerDef, err := t.ConvertDefDTOToTimerDefPB(defDTO)
		if err != nil {
			return nil, err
		}
		timerDefList = append(timerDefList, timerDef)
	}
	return timerDefList, nil
}
