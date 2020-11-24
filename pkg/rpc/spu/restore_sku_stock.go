package spu

import (
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	sku_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku"
)

func (controller *Controller) RestoreSkuStock(ctx context.Context, req *protos.RestoreSkuStockRequest) (*protos.RestoreSkuStockReply, error) {
	sku, err := sku_service.FindById(req.SkuId)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("get sku info failed")
		return &protos.RestoreSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	if sku == nil {
		logrus.WithField("skuId", req.SkuId).Warn("sku not found")
		return &protos.RestoreSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.EntityNotFound,
				ErrorMessageForDeveloper: errorcode.EntityNotFoundMsg,
				ErrorMessageForUser:      errorcode.EntityNotFoundMsg,
			},
		}, nil
	}

	newStock, err := sku_service.IncrStock(sku.SpuId, sku.ID, req.Count)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("restore sku stock failed")
		return &protos.RestoreSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.RestoreSkuStockReply{
		Stock: uint64(newStock),
	}, nil
}
