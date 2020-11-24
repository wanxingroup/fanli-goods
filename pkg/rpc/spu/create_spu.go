package spu

import (
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/restful/common/validator"

	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) CreateSPU(ctx context.Context, req *protos.CreateSPURequest) (*protos.CreateSPUReply, error) {

	var (
		err    error
		logger = log.GetLogger()
	)

	err = controller.ValidateCreateSPU(req)
	if err != nil {

		return &protos.CreateSPUReply{Err: controller.convertToProtoError(logger, err)}, nil
	}

	skus := make([]*product_model.Sku, 0, 10)

	// 目前仅一个sku
	sku := &product_model.Sku{
		ShopId:        req.GetShopId(),
		CategoryId:    req.GetCategoryId(),
		Price:         req.GetPrice(),
		VipPrice:      req.GetVipPrice(),
		Point:         req.GetPoint(),
		OriginalPrice: req.GetPrice(),
		Stock:         req.GetStock(),
		Name:          req.GetName(),
		Status:        product_model.Sku_Status_Normal,
	}
	sku.CreatedBy = req.GetCreatedBy()
	sku.UpdatedBy = req.GetCreatedBy()

	skus = append(skus, sku)

	spu := &product_model.Spu{
		ShopId:            req.GetShopId(),
		CategoryId:        req.GetCategoryId(),
		Name:              strings.TrimSpace(req.GetName()),
		BannerImages:      strings.TrimSpace(strings.Join(req.GetImages(), ",")),
		DetailText:        req.GetDetailText(),
		Price:             req.GetPrice(),
		VipPrice:          req.GetVipPrice(),
		SettlePrice:       req.GetSettlePrice(),
		CostPrice:         req.GetCostPrice(),
		IsHot:             req.GetIsHot(),
		Point:             req.GetPoint(),
		Stock:             req.GetStock(),
		Unit:              strings.TrimSpace(req.GetUnit()),
		Barcode:           strings.TrimSpace(req.GetBarcode()),
		Support7DayReturn: req.GetSupport7DayReturnGoods(),
		Status:            product_model.SPUStatusOffline,
		Skus:              skus,
		Sort:              uint(req.GetSort()),
	}
	spu.UpdatedBy = req.GetCreatedBy()
	spu.CreatedBy = req.GetCreatedBy()

	err = spu_service.Create(spu)
	if err != nil {
		logger.WithError(err).Error("create spu failed")
		return &protos.CreateSPUReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.CreateSpuError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.CreateSPUReply{
		Spu: controller.toSPUInformation(spu),
	}, nil
}

func (controller *Controller) ValidateCreateSPU(req *protos.CreateSPURequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ShopId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.CategoryId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.CategoryIdError), errorcode.CategoryIdErrorMessage)),
		),
		validation.Field(&req.Name, validation.Required, validator.NotBlank.
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUNameEmpty), errorcode.SPUNameEmptyMessage)),
		),
		validation.Field(&req.Price, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUPriceInvalid), errorcode.SPUPriceInvalidMessage)),
		),
		validation.Field(&req.Unit, validation.Required, validation.RuneLength(0, 4).
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUUnitTooLarge), errorcode.SPUUnitTooLargeMessage)),
		),
		validation.Field(&req.Images, validation.Required,
			validation.Length(0, 5).
				ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUImagesCountTooLarge), errorcode.SPUImagesCountTooLargeMessage)),
			validation.Length(1, 0).
				ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUImagesRequired), errorcode.SPUImagesRequiredMessage)),
		),
		validation.Field(&req.Stock, validation.Required, validation.Min(uint64(0)).
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUStockInvalid), errorcode.SPUStockInvalidMessage)),
		),
		validation.Field(&req.Sort,
			validation.Min(uint64(0)).Exclusive().ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUPriceInvalid), errorcode.SPUPriceInvalidMessage)),
			validation.Max(uint64(9999)).Exclusive().ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUPriceInvalid), errorcode.SPUPriceInvalidMessage)),
		),
	)
}
