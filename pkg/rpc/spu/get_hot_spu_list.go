package spu

import (
	"strings"

	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	productModel "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spuService "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) GetHotSPUList(ctx context.Context, req *protos.GetHotSPUListRequest) (*protos.GetHotSPUListReply, error) {

	logger := log.GetLogger().WithField("requestData", req)
	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
		"isHot":  1,
		"status": productModel.SPUStatusOnline,
	}

	pageData := map[string]uint64{
		"pageNum":  1,
		"pageSize": 10,
	}

	spus, _, err := spuService.GetSpuListByGetCondition(conditions, pageData)
	if err != nil {
		logger.WithError(err).Error("get spu infos failed")
		return &protos.GetHotSPUListReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	protoSpus := make([]*protos.SPUInformation, len(spus))
	for key, spu := range spus {
		protoSpus[key] = &protos.SPUInformation{
			SpuId:                  spu.ID,
			Name:                   spu.Name,
			CategoryId:             spu.CategoryId,
			Images:                 strings.Split(spu.BannerImages, ","),
			Price:                  spu.Price,
			VipPrice:               spu.VipPrice,
			Point:                  spu.Point,
			IsOnline:               spu.Status == product.SPUStatusOnline,
			Stock:                  spu.Stock,
			Unit:                   spu.Unit,
			ShopId:                 spu.ShopId,
			Support7DayReturnGoods: spu.Support7DayReturn,
			Barcode:                spu.Barcode,
			SettlePrice:            spu.SettlePrice,
			CostPrice:              spu.CostPrice,
			IsHot:                  spu.IsHot,
			SkuList:                []*protos.SKUInformationStruct{},
		}
	}

	controller.patchSKU(logger, protoSpus)

	return &protos.GetHotSPUListReply{
		Spus: protoSpus,
	}, nil
}
