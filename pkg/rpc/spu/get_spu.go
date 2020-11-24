package spu

import (
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) GetSpu(ctx context.Context, req *protos.GetSpuRequest) (*protos.GetSpuReply, error) {
	spu, err := spu_service.FindByShopIdAndSPUIdWithPreload(req.GetShopId(), req.GetSpuId())
	logger := log.GetLogger().WithField("req", req)

	if err != nil {
		logger.WithError(err).Error("get spu info failed")
		return &protos.GetSpuReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	if spu == nil {
		logger.Warn("spu not found")
		return &protos.GetSpuReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.EntityNotFound,
				ErrorMessageForDeveloper: errorcode.EntityNotFoundMsg,
				ErrorMessageForUser:      errorcode.EntityNotFoundMsg,
			},
		}, nil
	}

	return &protos.GetSpuReply{
		Spu: controller.toSPUInformation(spu),
	}, nil
}
