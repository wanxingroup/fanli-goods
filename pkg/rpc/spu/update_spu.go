package spu

import (
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"github.com/shomali11/util/xstrings"
	"golang.org/x/net/context"

	categoryPkg "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/restful/common/validator"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) UpdateSPU(ctx context.Context, req *protos.UpdateSPURequest) (*protos.UpdateSPUReply, error) {

	var (
		logger = log.GetLogger()
	)

	if err := controller.ValidateUpdateSPU(req); err != nil {

		return &protos.UpdateSPUReply{Err: controller.convertToProtoError(logger, err)}, nil
	}

	spu, replyError := controller.loadShopSPU(logger, req.GetShopId(), req.GetSpuId())
	if replyError != nil {
		return &protos.UpdateSPUReply{
			Err: replyError,
		}, nil
	}

	if spu == nil {

		return &protos.UpdateSPUReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.SpuNotFoundError,
				ErrorMessageForDeveloper: errorcode.SPUNotFoundErrorMessage,
				ErrorMessageForUser:      errorcode.SPUNotFoundErrorMessage,
			},
		}, nil
	}

	if xstrings.IsNotBlank(req.GetName()) {
		spu.Name = strings.TrimSpace(req.GetName())
	}
	if req.GetPrice() > 0 {
		spu.Price = req.GetPrice()
	}
	if req.GetVipPrice() >= 0 {
		spu.VipPrice = req.GetVipPrice()
	}
	if req.GetPoint() >= 0 {
		spu.Point = req.GetPoint()
	}
	if xstrings.IsNotBlank(req.GetUnit()) {
		spu.Unit = strings.TrimSpace(req.GetUnit())
	}
	if len(req.GetImages()) > 0 {
		spu.BannerImages = strings.Join(req.GetImages(), ",")
	}
	if xstrings.IsNotBlank(req.GetDetailText()) {
		spu.DetailText = req.GetDetailText()
	}

	spu.Support7DayReturn = req.GetSupport7DayReturnGoods()

	if xstrings.IsNotBlank(req.GetBarcode()) {
		spu.Barcode = strings.TrimSpace(req.GetBarcode())
	}

	spu.IsHot = req.GetIsHot()
	if req.GetSettlePrice() >= 0 {
		spu.SettlePrice = req.GetSettlePrice()
	}

	if req.GetCostPrice() >= 0 {
		spu.CostPrice = req.GetCostPrice()
	}
	spu.Sort = uint(req.GetSort())

	// 判断修改的SPU 分类是否存在
	if controller.CheckCategoryIdIsExist(req.GetCategoryId()) {
		spu.CategoryId = req.GetCategoryId()
	}

	// 当前只有一个sku
	skus, err := spu_service.FindSkus(req.GetSpuId())
	if err == nil {

		for _, sku := range skus {

			if req.GetPrice() > 0 {
				sku.Price = req.GetPrice()
				sku.OriginalPrice = req.GetPrice()
			}
			if req.GetVipPrice() >= 0 {
				sku.VipPrice = req.GetVipPrice()
			}
			if req.GetPoint() >= 0 {
				sku.Point = req.GetPoint()
			}
			if xstrings.IsNotBlank(req.GetName()) {
				sku.Name = strings.TrimSpace(req.GetName())
			}
		}

		spu.Skus = skus

	} else if !gorm.IsRecordNotFoundError(err) {

		logger.WithError(err).Error("get sku info failed")
		return &protos.UpdateSPUReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	spu.UpdatedBy = req.GetUpdatedBy()
	err = spu_service.Save(spu)
	if err != nil {

		logger.WithField("spuId", spu.ID).WithError(err).Error("update spu failed")
		return &protos.UpdateSPUReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.UpdateSpuError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.UpdateSpuErrorMessage,
			},
		}, nil
	}

	return &protos.UpdateSPUReply{
		Spu: controller.toSPUInformation(spu),
	}, nil
}

func (controller *Controller) ValidateUpdateSPU(req *protos.UpdateSPURequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.SpuId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.ShopId, validation.Required, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.ShopIdInvalid), errorcode.ShopIdInvalidMessage)),
		),
		validation.Field(&req.Name, validator.NotBlank.
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUNameEmpty), errorcode.SPUNameEmptyMessage)),
		),
		validation.Field(&req.Price, validation.Min(uint64(0)).Exclusive().
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUPriceInvalid), errorcode.SPUPriceInvalidMessage)),
		),
		validation.Field(&req.Images,
			validation.Length(0, 5).
				ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUImagesCountTooLarge), errorcode.SPUImagesCountTooLargeMessage)),
		),
		validation.Field(&req.Unit, validation.RuneLength(0, 4).
			ErrorObject(validation.NewError(strconv.Itoa(errorcode.SPUUnitTooLarge), errorcode.SPUUnitTooLargeMessage)),
		),
	)
}

func (controller *Controller) CheckCategoryIdIsExist(id uint64) bool {

	category, err := categoryPkg.FindById(id)

	if err == nil && category != nil {
		return true
	}
	return false
}
