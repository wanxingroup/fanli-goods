package spu

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	sku_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku"
	stock_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
)

func (controller *Controller) GetSkuStock(ctx context.Context, req *protos.GetSkuStockRequest) (*protos.GetSkuStockReply, error) {
	sku, err := sku_service.FindById(req.SkuId)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("get sku info failed")
		return &protos.GetSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	if sku == nil {
		logrus.WithField("skuId", req.SkuId).Warn("sku not found")
		return &protos.GetSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.EntityNotFound,
				ErrorMessageForDeveloper: errorcode.EntityNotFoundMsg,
				ErrorMessageForUser:      errorcode.EntityNotFoundMsg,
			},
		}, nil
	}

	stock, err := stock_service.GetRealStock(sku.SpuId, sku.ID)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("get sku stock failed")
		return &protos.GetSkuStockReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.GetStockError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.GetSkuStockReply{
		Stock: uint64(stock),
	}, nil
}
