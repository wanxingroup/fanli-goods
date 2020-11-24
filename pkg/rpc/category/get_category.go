package category

import (
	"github.com/jinzhu/gorm"
	context "golang.org/x/net/context"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (c *Controller) GetCategory(ctx context.Context, req *protos.GetCategoryRequest) (*protos.GetCategoryReply, error) {

	var (
		logger = log.GetLogger()
	)

	if req.GetShopId() <= 0 {

		logger.Info("shop id invalid")
		return &protos.GetCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "店铺 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	if req.GetCategoryId() <= 0 {

		logger.Info("category id invalid")
		return &protos.GetCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "分类 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	category, err := category_service.FindByShopIdAndCategoryId(req.GetShopId(), req.GetCategoryId())

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Info("category not exist")
			return &protos.GetCategoryReply{
				Err: &protos.Error{
					ErrorCode:                errorcode.CategoryNotFoundError,
					ErrorMessageForDeveloper: errorcode.CategoryNotFoundErrorMessage,
					ErrorMessageForUser:      errorcode.CategoryNotFoundErrorMessage,
				},
			}, nil
		}

		logger.WithError(err).Error("get category failed")
		return &protos.GetCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: errorcode.InternalErrorMsg,
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.GetCategoryReply{
		Category: c.toCategoryReplyObject(category),
	}, nil
}
