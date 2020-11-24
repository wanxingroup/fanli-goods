package category

import (
	"errors"
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	spu_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
)

func Create(category *product_model.Category) error {
	path := ","
	level := uint8(1)
	if category.ParentId > 0 {
		parentCategory, err := FindById(category.ParentId)
		if err != nil {
			return err
		}
		if parentCategory == nil {
			return errors.New("parent category not found")
		}
		path = fmt.Sprintf("%s%v,", parentCategory.Path, category.ParentId)
		level = parentCategory.Level + 1
	}
	category.Path = path
	category.Level = level
	category.ID = idcreator.NextID()
	if err := database.GetDB(constants.DatabaseConfigKey).Create(category).Error; err != nil {
		return err
	}
	return nil
}

func Save(category *product_model.Category) error {
	return database.GetDB(constants.DatabaseConfigKey).Save(category).Error
}

func UpdateStatus(categoryId uint64, status uint8) error {
	category := &product_model.Category{}
	category.ID = categoryId
	return database.GetDB(constants.DatabaseConfigKey).Model(category).Update("status", status).Error
}

func IsBindGoods(categoryId, shopId uint64) bool {

	spu, err := spu_service.FindByShopIdAndCategoryId(shopId, categoryId)

	if spu == nil && err == nil {
		return false
	}

	return true
}

func FindById(id uint64) (*product_model.Category, error) {
	category := new(product_model.Category)
	err := database.GetDB(constants.DatabaseConfigKey).First(category, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func FindByIds(ids []uint64) ([]*product_model.Category, error) {
	categories := make([]*product_model.Category, 0, len(ids))
	if len(ids) == 0 {
		return categories, nil
	}
	err := database.GetDB(constants.DatabaseConfigKey).Model(&product_model.Category{}).Find(&categories, "id in (?)", ids).Error
	return categories, err
}

// 特定查询的方法，方法名先定义为这个，后续尝试抽象为一个通用的条件查询
func FindByGetCondition(conditions map[string]interface{}) ([]*product_model.Category, error) {
	categories := make([]*product_model.Category, 0, 20)
	db := database.GetDB(constants.DatabaseConfigKey)
	if shopId, has := conditions["shopId"]; has {
		db = db.Where("shopId = ?", shopId)
	}
	if parentId, has := conditions["parentId"]; has {
		db = db.Where("parentId = ?", parentId)
	}
	if parentId, has := conditions["parentId"]; has {
		db = db.Where("parentId = ?", parentId)
	}
	db = db.Where("status = ?", product_model.Category_Status_Normal)
	err := db.Order("level, sort desc, id").Find(&categories).Error
	if gorm.IsRecordNotFoundError(err) {
		return categories, nil
	}

	return categories, err
}

func FindByShopIdAndCategoryId(shopId, categoryId uint64) (category *product_model.Category, err error) {

	category = &product_model.Category{}
	err = database.GetDB(constants.DatabaseConfigKey).
		Model(&product_model.Category{}).
		Where(&product_model.Category{
			BaseModel: databases.BaseModel{ID: categoryId},
			ShopId:    shopId,
		}).
		First(category).Error

	return
}
