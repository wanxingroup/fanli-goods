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

func (controller *Controller) GetSkus(ctx context.Context, req *protos.GetSkusRequest) (*protos.GetSkusReply, error) {

	if len(req.SkuIds) == 0 {
		return &protos.GetSkusReply{
			Skus: map[uint64]*protos.SKUInformation{},
		}, nil
	}

	skus, err := sku_service.FindByIds(req.SkuIds)

	if err != nil {
		logrus.WithField("skuIds", req.SkuIds).WithError(err).Error("get sku infos failed")
		return &protos.GetSkusReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	spuIds := make([]uint64, 0, len(skus))
	for _, sku := range skus {
		spuIds = append(spuIds, sku.SpuId)
	}

	spus, err := spu_service.FindByIds(spuIds)

	if err != nil {
		logrus.WithField("spuIds", spuIds).WithError(err).Error("get spu infos failed")
		return &protos.GetSkusReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	spuMap := make(map[uint64]*product.Spu, len(spus))
	for _, spu := range spus {
		spuMap[spu.ID] = spu
	}

	protoSkus := make(map[uint64]*protos.SKUInformation, len(skus))
	for _, sku := range skus {
		protoSkus[sku.ID] = &protos.SKUInformation{
			SkuId:         sku.ID,
			SpuId:         sku.SpuId,
			OriginalPrice: sku.OriginalPrice,
			Price:         sku.Price,
			VipPrice:      sku.VipPrice,
			Point:         sku.Point,
			Stock:         sku.Stock,
			Name:          sku.Name,
		}
		spu := spuMap[sku.SpuId]
		if spu != nil {
			protoSkus[sku.ID].Spu = &protos.SPUInformation{
				SpuId:      spu.ID,
				Name:       spu.Name,
				CategoryId: spu.CategoryId,
				Images:     strings.Split(spu.BannerImages, ","),
				Price:      uint64(spu.Price),
				VipPrice:   uint64(spu.VipPrice),
				Point:      uint64(spu.Point),
				IsOnline:   spu.Status == product.SPUStatusOnline,
			}
		}
	}
	return &protos.GetSkusReply{
		Skus: protoSkus,
	}, nil
}
