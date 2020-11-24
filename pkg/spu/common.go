package spu

import (
	"strings"

	"github.com/shomali11/util/xstrings"

	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
)

func ToSPUInformation(spu *product_model.Spu) *protos.SPUInformation {

	reply := &protos.SPUInformation{
		ShopId:                 spu.ShopId,
		SpuId:                  spu.ID,
		CategoryId:             spu.CategoryId,
		Name:                   spu.Name,
		Price:                  spu.Price,
		VipPrice:               spu.VipPrice,
		Point:                  spu.Point,
		Unit:                   spu.Unit,
		DetailText:             spu.DetailText,
		Stock:                  spu.Stock,
		Support7DayReturnGoods: spu.Support7DayReturn,
		Barcode:                spu.Barcode,
		IsOnline:               spu.Status == product_model.SPUStatusOnline,
		CreateTime:             spu.CreatedAt.Unix(),
		UpdateTime:             spu.UpdatedAt.Unix(),
	}

	if xstrings.IsNotEmpty(spu.BannerImages) {
		reply.Images = strings.Split(spu.BannerImages, ",")
	}

	if spu.Skus != nil {
		skuList := make([]*protos.SKUInformationStruct, 0, len(spu.Skus))
		for _, sku := range spu.Skus {
			skuList = append(skuList, ToSKUInformationStruct(sku))
		}
		reply.SkuList = skuList
	}
	return reply
}

func ToSKUInformationStruct(sku *product_model.Sku) *protos.SKUInformationStruct {

	return &protos.SKUInformationStruct{
		SkuId:         sku.ID,
		OriginalPrice: sku.OriginalPrice,
		Price:         sku.Price,
		VipPrice:      sku.VipPrice,
		Point:         sku.Point,
		Stock:         sku.Stock,
	}
}
