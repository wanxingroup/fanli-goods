package category

import (
	"strings"

	"github.com/shomali11/util/xstrings"
	context "golang.org/x/net/context"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (c *Controller) UpdateCategory(ctx context.Context, req *protos.UpdateCategoryRequest) (*protos.UpdateCategoryReply, error) {

	var (
		shopId, categoryId uint64
		logger             = log.GetLogger()
	)

	if req.GetShopId() <= 0 {

		logger.Info("shop id invalid")
		return &protos.UpdateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "店铺 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	if req.GetCategoryId() <= 0 {

		logger.Info("category id invalid")
		return &protos.UpdateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "分类 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	var category *product_model.Category
	shopId = req.GetShopId()
	categoryId = req.GetCategoryId()
	category, replyError := c.loadShopCategory(logger, shopId, categoryId)
	if replyError != nil {
		return &protos.UpdateCategoryReply{
			Err: replyError,
		}, nil
	}

	if category == nil {
		return &protos.UpdateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.CategoryNotFoundError,
				ErrorMessageForDeveloper: errorcode.CategoryNotFoundErrorMessage,
				ErrorMessageForUser:      errorcode.CategoryNotFoundErrorMessage,
			},
		}, nil
	}

	if xstrings.IsNotBlank(req.GetName()) {
		category.Name = strings.TrimSpace(req.GetName())
	}

	category.Description = strings.TrimSpace(req.GetDescription())

	if req.GetSort() >= 0 {
		category.Sort = uint16(req.GetSort())
	}

	category.UpdatedBy = req.GetUpdatedBy()
	err := category_service.Save(category)
	if err != nil {
		logger.WithField("categoryId", categoryId).WithError(err).Error("update category failed")
		return &protos.UpdateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.UpdateCategoryError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.UpdateCategoryErrorMessage,
			},
		}, nil
	}

	return &protos.UpdateCategoryReply{
		Category: c.toCategoryReplyObject(category),
	}, nil
}
