package category

import (
	"github.com/sirupsen/logrus"

	category_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/category"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/errorcode"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
)

func (c *Controller) toCategoryReplyObject(category *product_model.Category) *protos.CategoryInformation {
	return &protos.CategoryInformation{
		CategoryId:  category.ID,
		ShopId:      category.ShopId,
		ParentId:    category.ParentId,
		Name:        category.Name,
		Description: category.Description,
		Sort:        int64(category.Sort),
		Key:         category.Key,
		Level:       uint64(category.Level),
		Path:        category.Path,
		Unit:        category.Unit,
		Status:      uint64(category.Status),
	}
}

func (c *Controller) toCategoryReplyObjectList(categories []*product_model.Category) []*protos.CategoryInformation {
	categoryVos := make([]*protos.CategoryInformation, 0, len(categories))
	for _, category := range categories {
		categoryVos = append(categoryVos, c.toCategoryReplyObject(category))
	}
	return categoryVos
}

func (c *Controller) loadShopCategory(logger *logrus.Entry, shopId, categoryId uint64) (*product_model.Category, *protos.Error) {

	category, err := category_service.FindById(categoryId)
	if err != nil {
		logger.WithField("categoryId", categoryId).WithError(err).Info("get category info failed")
		return nil, &protos.Error{
			ErrorCode:                errorcode.CategoryNotFoundError,
			ErrorMessageForDeveloper: errorcode.CategoryNotFoundErrorMessage,
			ErrorMessageForUser:      errorcode.CategoryNotFoundErrorMessage,
		}
	}

	if category == nil || category.Status == product_model.Category_Status_Del {
		logger.WithField("category", category).Info("category was deleted")
		return nil, &protos.Error{
			ErrorCode:                errorcode.CategoryNotFoundError,
			ErrorMessageForDeveloper: "类目已被删除",
			ErrorMessageForUser:      errorcode.CategoryNotFoundErrorMessage,
		}
	}

	if category.ShopId != shopId {
		logger.WithField("category", category).WithField("shopId", shopId).Info("access forbidden")
		return nil, &protos.Error{
			ErrorCode:                errorcode.CategoryForbidError,
			ErrorMessageForDeveloper: "店铺 ID 跟类目店铺 ID 不一致",
			ErrorMessageForUser:      errorcode.CategoryForbidErrorMessage,
		}
	}
	return category, nil
}
