// Package rpc 负责应用对外的 rpc 通讯能力。
package rpc

import (
	"context"
	"fmt"

	"github.com/fflow-tech/fflow/service/cmd/foundation/timer/convertor"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	pb "github.com/fflow-tech/fflow/api/foundation/timer"
)

// TimerService 定时器后端服务实现
type TimerService struct {
	pb.UnimplementedEndpointServer
	domainService *service.DomainService
}

// NewTimerService 定时器构造函数
func NewTimerService(domainService *service.DomainService) *TimerService {
	return &TimerService{domainService: domainService}
}

// CreateTimer 创建定时器定义
func (w *TimerService) CreateTimer(ctx context.Context, req *pb.CreateTimerReq) (*pb.CreateTimerRsp, error) {
	rsp := &pb.CreateTimerRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	createDefDTO, err := convertor.TimerDefConvertor.ConvertCreatePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	defID, err := w.domainService.Commands.CreateTimerDef(createDefDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	rsp.Data = utils.Uint64ToStr(defID)
	return rsp, nil
}

// EnableTimer 激活定时器
func (w *TimerService) EnableTimer(ctx context.Context, req *pb.EnableTimerReq) (*pb.EnableTimerRsp, error) {
	rsp := &pb.EnableTimerRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	changeDTO, err := convertor.TimerDefConvertor.ConvertEnablePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := w.domainService.Commands.ChangeTimerStatus(changeDTO); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// DisableTimer 去激活定时器
func (w *TimerService) DisableTimer(ctx context.Context, req *pb.DisableTimerReq) (*pb.DisableTimerRsp, error) {
	rsp := &pb.DisableTimerRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	changeDTO, err := convertor.TimerDefConvertor.ConvertDisablePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := w.domainService.Commands.ChangeTimerStatus(changeDTO); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetTimerDef 获取定时器定义
func (w *TimerService) GetTimerDef(ctx context.Context, req *pb.GetTimerDefReq) (*pb.GetTimerDefRsp, error) {
	rsp := &pb.GetTimerDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getDTO, err := convertor.TimerDefConvertor.ConvertGetDefPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	def, err := w.domainService.Queries.GetTimerDef(getDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	timerDef, err := convertor.TimerDefConvertor.ConvertDefDTOToTimerDefPB(def)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	rsp.TimerDef = timerDef
	return rsp, nil
}

// DeleteTimerDef 删除定时器.
func (w *TimerService) DeleteTimerDef(ctx context.Context, req *pb.DeleteTimerDefReq) (*pb.DeleteTimerDefRsp, error) {
	rsp := &pb.DeleteTimerDefRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	delDTO, err := convertor.TimerDefConvertor.ConvertDelDefPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	if err := w.domainService.Commands.DeleteTimerDef(delDTO); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetTimerDefList 获取定时器定义列表
func (w *TimerService) GetTimerDefList(ctx context.Context, req *pb.GetTimerDefListReq) (
	*pb.GetTimerDefListRsp, error) {
	rsp := &pb.GetTimerDefListRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getDefListDTO, err := convertor.TimerDefConvertor.ConvertGetDefListPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	timerList, _, err := w.domainService.Queries.GetTimerDefList(getDefListDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	timerDefList, err := convertor.TimerDefConvertor.ConvertDefListDTOToTimerDefPB(timerList)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.TimerDefs = timerDefList
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetRunHistoryList 获取定时器执行列表
func (w *TimerService) GetRunHistoryList(ctx context.Context, req *pb.GetRunHistoryListReq) (
	*pb.GetRunHistoryListRsp, error) {
	rsp := &pb.GetRunHistoryListRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getHistoryListDTO, err := convertor.TimerDefConvertor.ConvertGetRunHistoryListPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	historyList, _, err := w.domainService.Queries.PageQueryHistory(getHistoryListDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	taskHistory, err := convertor.TimerDefConvertor.ConvertGetRunHistoryRspListDTOToPB(historyList)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.HistoryList = taskHistory
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// CreateApp 创建应用
func (w *TimerService) CreateApp(ctx context.Context, req *pb.CreateAppReq) (*pb.CreateAppRsp, error) {
	rsp := &pb.CreateAppRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	createAppDTO, err := convertor.AppConvertor.ConvertCreatePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	err = w.domainService.Commands.CreateApp(createAppDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// DeleteApp 删除应用
func (w *TimerService) DeleteApp(ctx context.Context, req *pb.DeleteAppReq) (*pb.DeleteAppRsp, error) {
	rsp := &pb.DeleteAppRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	deleteAppDTO, err := convertor.AppConvertor.ConvertDeletePbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	err = w.domainService.Commands.DeleteApp(deleteAppDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// GetAppList 获取应用列表
func (w *TimerService) GetAppList(ctx context.Context, req *pb.GetAppListReq) (*pb.GetAppListRsp, error) {
	rsp := &pb.GetAppListRsp{}
	if err := validateBasicReq(req.BasicReq); err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	getAppListDTO, err := convertor.AppConvertor.ConvertGetAppListPbToDTO(req)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.InvalidArgument.Code, err.Error())
		return rsp, nil
	}

	appList, _, err := w.domainService.Queries.GetAppList(getAppListDTO)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	appPBList, err := convertor.AppConvertor.ConvertAppDTOListToPB(appList)
	if err != nil {
		rsp.BasicRsp = NewFailedRsp(errno.Internal.Code, err.Error())
		return rsp, nil
	}

	rsp.AppList = appPBList
	rsp.BasicRsp = NewSucceedRsp()
	return rsp, nil
}

// NewSucceedRsp 生成成功返回
func NewSucceedRsp() *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
	}
}

// NewFailedRsp 通过自定义的错误码生成请求返回
func NewFailedRsp(code int32, message string) *pb.BasicRsp {
	return &pb.BasicRsp{
		Code:    code,
		Message: message,
	}
}

func validateBasicReq(req *pb.BasicReq) error {
	if req == nil || req.Namespace == "" {
		return fmt.Errorf("the req or req's namespace must not be empty")
	}

	return nil
}
