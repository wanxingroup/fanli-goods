package category

import (
	"strings"

	context "golang.org/x/net/context"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (c *Controller) CreateCategory(ctx context.Context, req *protos.CreateCategoryRequest) (*protos.CreateCategoryReply, error) {

	var (
		err    error
		logger = log.GetLogger()
	)

	if req.GetShopId() <= 0 {

		logger.Info("shop id invalid")
		return &protos.CreateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "店铺 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	if len(req.GetName()) <= 0 {
		logger.Info("category name empty")
		return &protos.CreateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "分类名称为空",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	category := &product_model.Category{
		ParentId:    req.GetParentId(),
		ShopId:      req.GetShopId(),
		Name:        strings.TrimSpace(req.GetName()),
		Description: strings.TrimSpace(req.GetDescription()),
		Sort:        uint16(req.GetSort()),
		Status:      product_model.Category_Status_Normal,
	}
	category.CreatedBy = req.GetCreatedBy()
	category.UpdatedBy = req.GetCreatedBy()
	err = category_service.Create(category)
	if err != nil {
		logger.WithError(err).Error("create category failed")
		return &protos.CreateCategoryReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.CreateCategoryError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.CreateCategoryErrorMessage,
			},
		}, nil
	}

	return &protos.CreateCategoryReply{
		Category: c.toCategoryReplyObject(category),
	}, nil
}
