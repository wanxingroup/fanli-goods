package spu

import (
	"strconv"
	"strings"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	productModel "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
)

func (controller *Controller) toSPUInformation(spu *productModel.Spu) *protos.SPUInformation {

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
		IsHot:                  spu.IsHot,
		SettlePrice:            spu.SettlePrice,
		CostPrice:              spu.CostPrice,
		IsOnline:               spu.Status == productModel.SPUStatusOnline,
		CreateTime:             spu.CreatedAt.Unix(),
		UpdateTime:             spu.UpdatedAt.Unix(),
		Sort:                   uint32(spu.Sort),
	}

	if xstrings.IsNotEmpty(spu.BannerImages) {
		reply.Images = strings.Split(spu.BannerImages, ",")
	}

	if spu.Skus != nil {
		skuList := make([]*protos.SKUInformationStruct, 0, len(spu.Skus))
		for _, sku := range spu.Skus {
			skuList = append(skuList, controller.toSKUInformationStruct(sku))
		}
		reply.SkuList = skuList
	}
	return reply
}

func (controller *Controller) toSKUInformationStruct(sku *productModel.Sku) *protos.SKUInformationStruct {

	return &protos.SKUInformationStruct{
		SkuId:         sku.ID,
		OriginalPrice: sku.OriginalPrice,
		Price:         sku.Price,
		VipPrice:      sku.VipPrice,
		Point:         sku.Point,
		Stock:         sku.Stock,
	}
}

func (controller *Controller) loadShopSPU(logger *logrus.Entry, shopId, spuId uint64) (*productModel.Spu, *protos.Error) {
	spu, err := spu_service.FindByIdWithPreload(spuId)
	if err != nil {
		logger.WithField("spuId", spuId).WithError(err).Error("get spu info failed")
		return nil, &protos.Error{
			ErrorCode:                errorcode.SpuNotFoundError,
			ErrorMessageForDeveloper: errorcode.SPUNotFoundErrorMessage,
			ErrorMessageForUser:      errorcode.SPUNotFoundErrorMessage,
		}
	}

	if spu == nil || spu.Status == productModel.SPUStatusDelete {
		logger.WithField("spuId", spuId).Info("spu no longer exist")
		return nil, &protos.Error{
			ErrorCode:                errorcode.SPUNoLongerExist,
			ErrorMessageForDeveloper: errorcode.SPUNoLongerExistMessage,
			ErrorMessageForUser:      errorcode.SPUNoLongerExistMessage,
		}
	}

	if spu.ShopId != shopId {
		logger.WithField("spuId", spuId).WithField("spu", spu).Info("access forbidden")
		return nil, &protos.Error{
			ErrorCode:                errorcode.SPUForbiddenError,
			ErrorMessageForDeveloper: errorcode.SPUForbiddenErrorMessage,
			ErrorMessageForUser:      errorcode.SPUForbiddenErrorMessage,
		}
	}
	return spu, nil
}

func (controller *Controller) convertToProtoError(logger *logrus.Entry, err error) *protos.Error {

	if err == nil {
		return nil
	}

	logger.WithError(err).Info("validate request data error")
	if validationError, ok := err.(validation.Error); ok {

		errCode, convertCodeError := strconv.Atoi(validationError.Code())
		if convertCodeError != nil {
			logger.WithError(convertCodeError).Error("convert validation error code error")
			errCode = errorcode.RequestParameterError
		}

		return &protos.Error{
			ErrorCode:                uint32(errCode),
			ErrorMessageForDeveloper: validationError.Message(),
			ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
		}
	}

	return &protos.Error{
		ErrorCode:                errorcode.RequestParameterError,
		ErrorMessageForDeveloper: err.Error(),
		ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
	}
}

func (controller *Controller) patchSKU(logger *logrus.Entry, spus []*protos.SPUInformation) {

	if len(spus) == 0 {
		return
	}
	spuIds := make([]uint64, 0, len(spus))

	checkIds := map[uint64]*protos.SPUInformation{}

	for _, spu := range spus {
		spuIds = append(spuIds, spu.SpuId)
		checkIds[spu.SpuId] = spu
	}

	skus := []*productModel.Sku{}
	err := database.GetDB(constants.DatabaseConfigKey).
		Model(&productModel.Sku{}).
		Where("`spuId` IN (?)", spuIds).
		Where(&productModel.Sku{Status: productModel.Sku_Status_Normal}).
		Find(&skus).
		Error

	if err != nil {
		logger.WithError(err).Error("get skus error")
		return
	}

	for _, sku := range skus {

		spu, exist := checkIds[sku.SpuId]
		if !exist {
			continue
		}

		if spu.SkuList == nil {
			spu.SkuList = make([]*protos.SKUInformationStruct, 0)
		}
		spu.SkuList = append(spu.SkuList, &protos.SKUInformationStruct{
			SkuId:         sku.ID,
			OriginalPrice: sku.OriginalPrice,
			Price:         sku.Price,
			VipPrice:      sku.VipPrice,
			Point:         sku.Point,
			Stock:         sku.Stock,
			Name:          sku.Name,
		})
	}
}
