package spu

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	sku_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
)

func (controller *Controller) GetSku(ctx context.Context, req *protos.GetSkuRequest) (*protos.GetSkuReply, error) {
	sku, err := sku_service.FindById(req.SkuId)

	if err != nil {
		logrus.WithField("skuId", req.SkuId).WithError(err).Error("get sku info failed")
		return &protos.GetSkuReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	if sku == nil {
		logrus.WithField("skuId", req.SkuId).Warn("sku not found")
		return &protos.GetSkuReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.EntityNotFound,
				ErrorMessageForDeveloper: errorcode.EntityNotFoundMsg,
				ErrorMessageForUser:      errorcode.EntityNotFoundMsg,
			},
		}, nil
	}

	spu, err := spu_service.FindById(sku.SpuId)

	if err != nil {
		logrus.WithField("spuId", sku.SpuId).WithError(err).Error("get spu info failed")
		return &protos.GetSkuReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.GetSkuReply{
		Sku: &protos.SKUInformation{
			SkuId:         sku.ID,
			SpuId:         sku.SpuId,
			OriginalPrice: sku.OriginalPrice,
			Price:         sku.Price,
			VipPrice:      sku.VipPrice,
			Point:         sku.Point,
			Stock:         sku.Stock,
			Name:          sku.Name,
			Spu: &protos.SPUInformation{
				SpuId:      spu.ID,
				Name:       spu.Name,
				CategoryId: spu.CategoryId,
				Images:     strings.Split(spu.BannerImages, ","),
				Price:      spu.Price,
				IsOnline:   spu.Status == product.SPUStatusOnline,
			},
		},
	}, nil
}
