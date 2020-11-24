package spu

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shomali11/util/xstrings"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

func (controller *Controller) GetSPUList(ctx context.Context, req *protos.GetSPUListRequest) (*protos.GetSPUListReply, error) {

	logger := log.GetLogger().WithField("requestData", req)

	var pageNum = uint64(defaultPageNum)
	var pageSize = uint64(defaultPageSize)

	if req.GetPage() > 0 {
		pageNum = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
	}

	if !validation.IsEmpty(req.GetStatus()) {
		conditions["status"] = req.GetStatus()
	}

	if !validation.IsEmpty(req.GetIsHot()) {
		if req.GetIsHot() {
			conditions["isHot"] = 1
		} else {
			conditions["isHot"] = 0
		}
	}

	if xstrings.IsNotBlank(req.GetNameFuzzySearch()) {
		conditions["nameFuzzySearch"] = strings.TrimSpace(req.GetNameFuzzySearch())
	}

	pageData := map[string]uint64{
		"page":     pageNum,
		"pageSize": pageSize,
	}

	if !validation.IsEmpty(req.GetLastSpuId()) {
		conditions["lastSpuId"] = req.GetLastSpuId()
	}

	spus, count, err := spu_service.GetSpuListByGetCondition(conditions, pageData)

	if err != nil {
		logger.WithError(err).Error("get spu infos failed")
		return &protos.GetSPUListReply{
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

	return &protos.GetSPUListReply{
		Spus:  protoSpus,
		Count: count,
	}, nil
}
