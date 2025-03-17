package web

import (
	"github.com/gin-gonic/gin"
	"github.com/fflow-tech/fflow/service/internal/demo-app/blank-demo/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/demo-app/blank-demo/domain/service"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/errno"
	"net/http"
)

// CrawlerController blank-demo-app/blank-demo http 服务实现
type CrawlerController struct {
	domainService *service.DomainService
}

// NewCrawlerController 构造函数
func NewCrawlerController(domainService *service.DomainService) *CrawlerController {
	return &CrawlerController{domainService: domainService}
}

// StartCollect 开始采集
// @Summary 采集相关接口
// @Description 采集相关接口
// @Tags 采集相关接口
// @Accept application/json
// @Produce application/json
// @Param callReq body dto.StartCollectReqDTO true "采集请求"
// @Success 200 {object} dto.StartCollectRspDTO
// @Router /blog-blank-demo-app/api/v1/collect/start [post]
func (h *CrawlerController) StartCollect(c *gin.Context) {
	var req dto.StartCollectReqDTO
	if err := bindReq(c, &req); err != nil {
		c.JSON(http.StatusOK, constants.NewFailedWebRspWithMsg(errno.InvalidArgument, err.Error()))
		return
	}
	err := h.domainService.Commands.StartCollect(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, constants.NewFailedWebRspWithMsg(errno.Unauthenticated, err.Error()))
		return
	}

	c.JSON(http.StatusOK, constants.NewSucceedWebRsp(nil))
}

func bindReq(c *gin.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	return nil
}
