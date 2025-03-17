package convertor

import (
	pb "github.com/fflow-tech/fflow/api/foundation/timer"
	"github.com/fflow-tech/fflow/service/internal/foundation/timer/domain/dto"
	"github.com/fflow-tech/fflow/service/pkg/constants"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"github.com/jinzhu/copier"
)

var (
	// AppConvertor 应用转换器
	AppConvertor = &appConvertor{}
)

type appConvertor struct {
}

// ConvertCreatePbToDTO 创建PB转换DTO
func (*appConvertor) ConvertCreatePbToDTO(req *pb.CreateAppReq) (*dto.CreateAppDTO, error) {
	createAppDTO := &dto.CreateAppDTO{}
	err := copier.Copy(createAppDTO, req)
	if err != nil {
		log.Errorf("Failed to APP ConvertCreatePbToDTO copy error, caused by %s, req:%s", err,
			utils.StructToJsonStr(req))
		return nil, err
	}
	return createAppDTO, nil
}

// ConvertDeletePbToDTO 删除PB转换DTO
func (*appConvertor) ConvertDeletePbToDTO(req *pb.DeleteAppReq) (*dto.DeleteAppDTO, error) {
	deleteAppDTO := &dto.DeleteAppDTO{}
	err := copier.Copy(deleteAppDTO, req)
	if err != nil {
		log.Errorf("Failed to APP ConvertDeletePbToDTO copy error, caused by %s, req:%s", err,
			utils.StructToJsonStr(req))
		return nil, err
	}
	return deleteAppDTO, nil
}

// ConvertGetAppListPbToDTO 获取应用列表PB转换DTO
func (*appConvertor) ConvertGetAppListPbToDTO(req *pb.GetAppListReq) (*dto.PageQueryAppDTO,
	error) {
	getAppListDTO := &dto.PageQueryAppDTO{}
	getAppListDTO.PageQuery = constants.NewPageQuery(int(req.PageIndex), int(req.PageSize))
	getAppListDTO.Order = constants.NewDefaultOrder()
	if req.AscOrder {
		getAppListDTO.Order.Order = constants.DescOrder
	}
	getAppListDTO.Name = req.Name
	getAppListDTO.Creator = req.Creator
	return getAppListDTO, nil
}

// ConvertAppDTOToPB 获取应用信息DTO转换PB
func (*appConvertor) ConvertAppDTOToPB(appDTO *dto.App) (*pb.AppInfo, error) {
	appInfo := &pb.AppInfo{}
	err := copier.Copy(appInfo, appDTO)
	if err != nil {
		log.Errorf("Failed to copy error, caused by %s, req:%s", err, utils.StructToJsonStr(appInfo))
		return nil, err
	}

	return appInfo, nil
}

// ConvertAppDTOListToPB 获取应用信息列表DTO转换PB
func (a *appConvertor) ConvertAppDTOListToPB(appDTOList []*dto.App) ([]*pb.AppInfo, error) {
	appInfoList := make([]*pb.AppInfo, 0, len(appDTOList))
	for _, appDTO := range appDTOList {
		appInfo, err := a.ConvertAppDTOToPB(appDTO)
		if err != nil {
			return nil, err
		}
		appInfoList = append(appInfoList, appInfo)
	}
	return appInfoList, nil
}
