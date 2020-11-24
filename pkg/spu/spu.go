package spu

import (
	"errors"
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/jinzhu/gorm"

	baseDb "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	productModel "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	stockService "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
)

var ErrorInvalidSpuStatus = errors.New("非法的商品状态")
var ErrorSpuHasNoSkus = errors.New("商品下无SKU")

func Create(spu *productModel.Spu) error {
	if spu.ID == 0 {
		spu.ID = idcreator.NextID()
	}
	skuStocks := make(map[uint64]uint64, len(spu.Skus))
	for _, sku := range spu.Skus {
		if sku.ID == 0 {
			sku.ID = idcreator.NextID()
		}
		skuStocks[sku.ID] = sku.Stock
	}
	// 级联创建skus
	err := database.GetDB(constants.DatabaseConfigKey).Create(spu).Error
	if err == nil {
		// 设置库存
		err = stockService.SetRealStocks(spu.ID, skuStocks)
	}
	return err
}

func Save(spu *productModel.Spu) error {
	tx := database.GetDB(constants.DatabaseConfigKey).Begin()

	if err := tx.Error; err != nil {
		return err
	}

	// 避免更新stock
	err := tx.Set("gorm:association_autoupdate", false).Model(spu).Omit("stock", "status").Save(spu).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if len(spu.Skus) > 0 {
		for _, sku := range spu.Skus {
			// 避免更新stock
			err := tx.Model(sku).Omit("stock", "status").Save(sku).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func FindById(id uint64) (*productModel.Spu, error) {
	spu := new(productModel.Spu)
	err := database.GetDB(constants.DatabaseConfigKey).First(spu, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return spu, nil
}

func FindByIdWithPreload(id uint64) (*productModel.Spu, error) {
	spu := new(productModel.Spu)
	err := database.GetDB(constants.DatabaseConfigKey).Preload("Skus").First(spu, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return spu, nil
}

func FindByShopIdAndSPUIdWithPreload(shopId, spuId uint64) (*productModel.Spu, error) {
	spu := new(productModel.Spu)
	err := database.GetDB(constants.DatabaseConfigKey).Preload("Skus").Where("id = ? AND shopId = ?", spuId, shopId).First(spu).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return spu, nil
}

func FindByShopIdAndCategoryId(shopId, categoryId uint64) (*productModel.Spu, error) {
	spu := new(productModel.Spu)
	err := database.GetDB(constants.DatabaseConfigKey).Preload("Skus").Where("categoryId = ? AND shopId = ? AND status != ?", categoryId, shopId, productModel.SPUStatusDelete).First(spu).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return spu, nil
}

func FindByIds(ids []uint64) ([]*productModel.Spu, error) {
	spus := make([]*productModel.Spu, 0, len(ids))
	if len(ids) == 0 {
		return spus, nil
	}
	err := database.GetDB(constants.DatabaseConfigKey).Model(&productModel.Spu{}).Find(&spus, "id in (?)", ids).Error
	return spus, err
}

func FindByShopIdAndSPUIds(shopId uint64, spuIds []uint64) ([]*productModel.Spu, error) {
	spus := make([]*productModel.Spu, 0, len(spuIds))
	if len(spuIds) == 0 {
		return spus, nil
	}
	err := database.GetDB(constants.DatabaseConfigKey).Model(&productModel.Spu{}).
		Find(&spus, "id in (?) and shopId = ?", spuIds, shopId).Error
	return spus, err
}

func FindSkus(spuId uint64) ([]*productModel.Sku, error) {
	skus := make([]*productModel.Sku, 0, 20)
	err := database.GetDB(constants.DatabaseConfigKey).Find(&skus, "spuId = ? and status = ?", spuId, productModel.Sku_Status_Normal).Error
	return skus, err
}

func UpdateStatus(spuId uint64, status uint8) error {
	if !IsValidStatus(status) {
		return ErrorInvalidSpuStatus
	}
	spu := &productModel.Spu{}
	spu.ID = spuId
	return database.GetDB(constants.DatabaseConfigKey).Model(spu).Update("status", status).Error
}

func SetStock(spuId uint64, stock uint64) error {
	// 针对只有一个sku的实现, 目前不保证db与redis的一致性, 简单实现
	tx := database.GetDB(constants.DatabaseConfigKey).Begin()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&productModel.Spu{}).Where("id = ?", spuId).Update("stock", stock).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	skus, err := FindSkus(spuId)

	if err != nil {
		tx.Rollback()
		return err
	}

	if len(skus) == 0 {
		tx.Rollback()
		return ErrorSpuHasNoSkus
	}

	err = database.GetDB(constants.DatabaseConfigKey).Model(skus[0]).Update("stock", stock).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	// 目前是在spu上设置库存，同步给仅有一个sku的情况，如果有多sku，那么需要从sku维度设置库存，然后取和设置给spu
	return stockService.SetRealStock(spuId, skus[0].ID, stock)
}

// 查询如果仅仅是等于，那么可以采用find by example，但会有其他的条件，所以是否可以封装个通用的

func GetSpuListByGetCondition(conditions map[string]interface{}, pageData map[string]uint64) ([]*productModel.Spu, uint64, error) {
	db := database.GetDB(constants.DatabaseConfigKey).Model(&productModel.Spu{})
	if shopId, has := conditions["shopId"]; has {
		db = db.Where("shopId = ?", shopId)
	}

	if categoryId, has := conditions["categoryId"]; has {
		db = db.Where("categoryId = ?", categoryId)
	}

	if lastSpuId, has := conditions["lastSpuId"]; has {
		db = db.Where("id < ?", lastSpuId)
	}

	if nameFuzzySearch, has := conditions["nameFuzzySearch"]; has {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", nameFuzzySearch))
	}

	if status, has := conditions["status"]; has {
		db = db.Where("status in (?)", status)
	}

	if isHot, has := conditions["isHot"]; has {
		db = db.Where("isHot = ?", isHot)
	}

	db = db.Order("sort")

	results := make([]*productModel.Spu, 0, pageData["pageSize"])
	var count uint64
	err := baseDb.FindPage(db, pageData, &results, &count)
	return results, count, err

}

// 这里先丑陋点，后面再重构
func FindByXcxCondition(conditions map[string]interface{}, fromSort uint, fromId uint64, limit uint) ([]*productModel.Spu, error) {
	results := make([]*productModel.Spu, 0, limit)
	db := database.GetDB(constants.DatabaseConfigKey).Preload("Skus", "status = ?", productModel.Sku_Status_Normal)
	if shopIds, has := conditions["shopIds"]; has && len(shopIds.([]uint64)) > 0 {
		db = db.Where("shopId in (?)", shopIds)
	}
	if categoryId, has := conditions["categoryId"]; has && categoryId.(uint64) > 0 {
		db = db.Where("categoryId = ?", categoryId)
	}
	if fromId > 0 {
		db = db.Where("sort < ? or (sort = ? and id > ?)", fromSort, fromSort, fromId)
	}
	db = db.Where("status = ?", productModel.SPUStatusOnline)
	err := db.Order("sort desc, id").Limit(limit).Find(&results).Error
	return results, err
}

func IsValidStatus(status uint8) bool {
	return productModel.SPUStatusDelete == status ||
		productModel.SPUStatusOffline == status ||
		productModel.SPUStatusOnline == status
}
