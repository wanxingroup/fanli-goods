package spu

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (controller *Controller) GetSpus(ctx context.Context, req *protos.GetSpusRequest) (*protos.GetSpusReply, error) {

	logger := log.GetLogger().WithField("requestData", req)

	if len(req.GetSpuIds()) == 0 {
		logger.Info("spu")
		return &protos.GetSpusReply{
			Spus: map[uint64]*protos.SPUInformation{},
		}, nil
	}

	spus, err := spu_service.FindByShopIdAndSPUIds(req.GetShopId(), req.GetSpuIds())

	if err != nil {
		logrus.WithField("requestData", req).WithError(err).Error("get spu infos failed")
		return &protos.GetSpusReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	protoSpus := make(map[uint64]*protos.SPUInformation, len(spus))
	for _, spu := range spus {
		protoSpus[spu.ID] = &protos.SPUInformation{
			SpuId:       spu.ID,
			Name:        spu.Name,
			CategoryId:  spu.CategoryId,
			Images:      strings.Split(spu.BannerImages, ","),
			Price:       spu.Price,
			VipPrice:    spu.VipPrice,
			Point:       spu.Point,
			IsHot:       spu.IsHot,
			SettlePrice: spu.SettlePrice,
			CostPrice:   spu.CostPrice,
			IsOnline:    spu.Status == product.SPUStatusOnline,
		}
	}
	return &protos.GetSpusReply{
		Spus: protoSpus,
	}, nil
}
