package spu

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) UpdateStock(ctx context.Context, req *protos.UpdateStockRequest) (*protos.UpdateStockReply, error) {

	var (
		logger = log.GetLogger()
	)

	if err := controller.ValidateUpdateStock(req); err != nil {

		return &protos.UpdateStockReply{Err: controller.convertToProtoError(logger, err)}, nil
	}

	spu, replyError := controller.loadShopSPU(logger, req.GetShopId(), req.GetSpuId())
	if replyError != nil {
		return &protos.UpdateStockReply{Err: replyError}, nil
	}

	if spu == nil {

		return &protos.UpdateStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.SpuNotFoundError,
				ErrorMessageForDeveloper: errorcode.SPUNotFoundErrorMessage,
				ErrorMessageForUser:      errorcode.SPUNotFoundErrorMessage,
			},
		}, nil
	}

	err := spu_service.SetStock(req.GetSpuId(), req.GetStock())
	if err != nil {
		logger.WithField("req", req).WithError(err).Error("set spu stock failed")
		return &protos.UpdateStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.UpdateStockReply{
		Stock: req.GetStock(),
	}, nil
}

func (controller *Controller) ValidateUpdateStock(req *protos.UpdateStockRequest) error {

	return validation.ValidateStruct(req,
		validation.Field(&req.SpuId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.ShopId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.Stock, validation.Required, validation.Min(uint64(0))),
	)
}
