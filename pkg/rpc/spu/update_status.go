package spu

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) UpdateStatus(ctx context.Context, req *protos.UpdateStatusRequest) (*protos.UpdateStatusReply, error) {

	var (
		logger = log.GetLogger()
	)

	logger.WithField("requestData", req).Info("start to update status")
	if err := controller.ValidateUpdateStatus(req); err != nil {

		return &protos.UpdateStatusReply{Err: controller.convertToProtoError(logger, err)}, nil
	}

	spu, replyError := controller.loadShopSPU(logger, req.GetShopId(), req.GetSpuId())
	if replyError != nil {
		return &protos.UpdateStatusReply{Err: replyError}, nil
	}

	if spu == nil {

		return &protos.UpdateStatusReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.SpuNotFoundError,
				ErrorMessageForDeveloper: errorcode.SPUNotFoundErrorMessage,
				ErrorMessageForUser:      errorcode.SPUNotFoundErrorMessage,
			},
		}, nil
	}

	spu.Status = uint8(req.GetStatus())
	err := spu_service.UpdateStatus(req.GetSpuId(), spu.Status)
	if err != nil {
		logger.WithField("req", req).WithError(err).
			Error("update spu status failed")
		return &protos.UpdateStatusReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.CreateSpuError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.UpdateStatusReply{}, nil
}

func (controller *Controller) ValidateUpdateStatus(req *protos.UpdateStatusRequest) error {

	return validation.ValidateStruct(req,
		validation.Field(&req.SpuId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.ShopId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.Status, validation.Required,
			validation.In(
				uint64(product_model.SPUStatusOnline), uint64(product_model.SPUStatusOffline),
			),
		),
	)
}
