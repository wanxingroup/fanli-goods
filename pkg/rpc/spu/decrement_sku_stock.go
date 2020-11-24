package spu

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	sku_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku"
	stock_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
)

func (controller *Controller) DecrementSkuStock(ctx context.Context, req *protos.DecrementSkuStockRequest) (*protos.DecrementSkuStockReply, error) {
	sku, err := sku_service.FindById(req.SkuId)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("get sku info failed")
		return &protos.DecrementSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	if sku == nil {
		logrus.WithField("skuId", req.SkuId).Warn("sku not found")
		return &protos.DecrementSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.EntityNotFound,
				ErrorMessageForDeveloper: errorcode.EntityNotFoundMsg,
				ErrorMessageForUser:      errorcode.EntityNotFoundMsg,
			},
		}, nil
	}

	newStock, err := sku_service.DecrStock(sku.SpuId, sku.ID, req.Count)

	if err != nil {
		if err, ok := err.(*stock_service.OutOfStockError); ok {
			return &protos.DecrementSkuStockReply{
				Err: &protos.Error{
					ErrorCode:                errorcode.OutOfStockError,
					ErrorMessageForDeveloper: err.Error(),
					ErrorMessageForUser:      errorcode.OutOfStockErrorMsg,
				},
			}, nil
		}

		logrus.WithField("skuId", req.SkuId).WithError(err).Error("decrement sku stock failed")
		return &protos.DecrementSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.DecrementSkuStockReply{
		Stock: uint64(newStock),
	}, nil
}
