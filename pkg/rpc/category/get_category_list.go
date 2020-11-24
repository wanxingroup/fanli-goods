package category

import (
	"golang.org/x/net/context"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func (c *Controller) GetCategoryList(ctx context.Context, req *protos.GetCategoryListRequest) (*protos.GetCategoryListReply, error) {

	// 由于每层分类不多，先不做分页
	var (
		logger = log.GetLogger()
	)

	if req.GetShopId() <= 0 {

		logger.Info("shop id invalid")
		return &protos.GetCategoryListReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.RequestParameterError,
				ErrorMessageForDeveloper: "店铺 ID 无效",
				ErrorMessageForUser:      errorcode.RequestParameterErrorMessage,
			},
		}, nil
	}

	conditions := map[string]interface{}{
		"shopId": req.GetShopId(),
	}
	if req.GetParentCategoryId() >= 0 { // parentId < 0 to get all
		conditions["parentId"] = req.GetParentCategoryId()
	}

	categories, err := category_service.FindByGetCondition(conditions)

	if err != nil {
		logger.WithError(err).Error("get category list failed")
		return &protos.GetCategoryListReply{
			Err: &protos.Error{
				ErrorCode:                errorcode.InternalError,
				ErrorMessageForDeveloper: err.Error(),
				ErrorMessageForUser:      errorcode.InternalErrorMsg,
			},
		}, nil
	}

	return &protos.GetCategoryListReply{
		Category: c.toCategoryReplyObjectList(categories),
	}, nil
}
