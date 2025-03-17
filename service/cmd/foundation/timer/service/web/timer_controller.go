package web

import (
	"net/http"
	"reflect"

	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TimerController workflow-app/engine http服务实现
type TimerController struct {
	domainService *service.DomainService
}

// NewTimerController 构造函数
func NewTimerController(domainService *service.DomainService) *TimerController {
	return &TimerController{domainService: domainService}
}

// CreateTimerDef 创建定时器定义
// @Summary 创建定时器定义
// @Description 创建定时器定义
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.CreateTimerDefDTO true "创建定时器定义"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/create [post]
func (h *TimerController) CreateTimerDef(c *gin.Context) {
	var req dto.CreateTimerDefDTO
	// 绑定参数
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Commands.CreateTimerDef(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}
	rspData := utils.Uint64ToStr(data)
	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(rspData))
}

// ChangeDefStatus 改变定时器状态
// @Summary 改变定时器状态
// @Description 改变定时器状态
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.ChangeTimerStatusDTO true "修改定时器状态"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/change [post]
func (h *TimerController) ChangeDefStatus(c *gin.Context) {
	var req dto.ChangeTimerStatusDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	err := h.domainService.Commands.ChangeTimerStatus(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(nil))
}

// GetDefDetail 查询单条定时器定义
// @Summary 查询单条定时器定义
// @Description 查询单条定时器定义
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.GetTimerDefDTO true "查询单条定时器定义"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/get [get]
func (h *TimerController) GetDefDetail(c *gin.Context) {
	var req dto.GetTimerDefDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	data, err := h.domainService.Queries.GetTimerDef(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(data))
}

// GetTimerDefList 分页查询定时器列表
// @Summary 分页查询定时器列表
// @Description 分页查询定时器列表
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.PageQueryTimeDefDTO true "分页查询定时器列表"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/list [get]
func (h *TimerController) GetTimerDefList(c *gin.Context) {
	var req dto.PageQueryTimeDefDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	timerDefList, total, err := h.domainService.Queries.GetTimerDefList(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(timerDefList, total))
}

// DeleteTimer 删除定时器
// @Summary 删除定时器
// @Description 删除定时器
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.DeleteTimerDefDTO true "删除定时器"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/delete [delete]
func (h *TimerController) DeleteTimer(c *gin.Context) {
	var req dto.DeleteTimerDefDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	if err := h.domainService.Commands.DeleteTimerDef(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(nil))
}

// GetTimerRunHistory 获取定时器运行历史记录
// @Summary 获取定时器运行历史记录
// @Description 获取定时器运行历史记录
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.PageQueryTimeDefDTO true "获取定时器运行历史记录"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/runHistory [GET]
func (h *TimerController) GetTimerRunHistory(c *gin.Context) {
	var req dto.PageQueryRunHistoryDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	timerRunHistory, total, err := h.domainService.Queries.PageQueryHistory(&req)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRspWithTotal(timerRunHistory, total))
}

// GetTimerTaskList 获取定时器任务列表
// @Summary 获取定时器任务列表
// @Description 获取定时器任务列表
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def query dto.GetTimerTaskListDTO true "获取定时器任务列表"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/timerTaskList [GET]
func (h *TimerController) GetTimerTaskList(c *gin.Context) {
	var req dto.GetTimerTaskListDTO
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	timerTaskList, err := h.domainService.Queries.GetTimeLimitTimers(req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(timerTaskList))
}

// TimerListSendNotify 定时器列表批量触发  此接口只能内部使用 谨慎使用可能会引起定时器重复触发的问题
// @Summary 定时器列表批量触发
// @Description 定时器列表批量触发
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Param def body dto.TimerListSendNotifyDTO true "定时器列表批量触发"
// @Success 200 {object} WebRsp
// @Router /timer/api/v1/def/timerListSend [post]
func (h *TimerController) TimerListSendNotify(c *gin.Context) {
	var req dto.TimerListSendNotifyDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}

	err := h.domainService.Commands.ManualTriggerSendList(req.TimerList)
	if err != nil {
		c.JSON(http.StatusOK, NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewSucceedWebRsp(nil))
}

// DeleteRunHistories 删除过期历史执行记录
// @Summary 删除过期历史执行记录
// @Description 删除过期历史执行记录
// @Tags 定时器相关接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} constants.WebRsp
// @Router /timer/api/v1/def/deleteRunHistories [delete]
func (h *TimerController) DeleteRunHistories(c *gin.Context) {
	err := h.domainService.Commands.DeleteRunHistories()
	if err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.Internal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp("delete run histories success"))
}

func bindReq(c *gin.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	return nil
}

func setOperator(req interface{}, name string) {
	v := reflect.ValueOf(req).Elem()
	if !v.IsValid() {
		return
	}

	creator := v.FieldByName("Creator")
	if creator.IsValid() {
		creator.Set(reflect.ValueOf(name))
	}

	operator := v.FieldByName("Operator")
	if operator.IsValid() {
		operator.Set(reflect.ValueOf(name))
	}
}
