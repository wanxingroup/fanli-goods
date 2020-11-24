package category

import (
	context "golang.org/x/net/context"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (c *Controller) DeleteCategory(ctx context.Context, req *protos.DeleteCategoryRequest) (*protos.DeleteCategoryReply, error) {

	var (
		logger = log.GetLogger()
	)

	if req.GetShopId() <= 0 {

		logger.Info("shop id invalid")
		return &protos.DeleteCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "店铺 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	if req.GetCategoryId() <= 0 {

		logger.Info("category id invalid")
		return &protos.DeleteCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "分类 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	category, errReply := c.loadShopCategory(logger, req.GetShopId(), req.GetCategoryId())

	if errReply != nil {
		return &protos.DeleteCategoryReply{
			Err: errReply,
		}, nil
	}

	if category == nil {

		logger.Info("category not exist")
		return &protos.DeleteCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.CategoryNotFoundError,
				ErrorMessageForDeveloper: "分类 ID 无效",
				ErrorMessageForUser:      errorcode.CategoryNotFoundErrorMessage,
			},
		}, nil
	}

	// 校验是否有商品绑定分类
	isBind := category_service.IsBindGoods(req.GetCategoryId(), req.GetShopId())
	if isBind {
		return &protos.DeleteCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.DeleteNotNilCategoryError,
				ErrorMessageForDeveloper: errorcode.DeleteNotNilCategoryErrorMessage,
				ErrorMessageForUser:      errorcode.DeleteNotNilCategoryErrorMessage,
			},
		}, nil
	}

	err := category_service.UpdateStatus(req.GetCategoryId(), product_model.Category_Status_Del)

	if err != nil {
		logger.WithError(err).Error("delete category failed")
		return &protos.DeleteCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.DeleteCategoryError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.DeleteCategoryErrorMessage,
			},
		}, nil
	}

	return &protos.DeleteCategoryReply{}, nil
}
