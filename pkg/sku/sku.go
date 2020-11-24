package sku

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
)

func FindById(id uint64) (*product_model.Sku, error) {
	sku := new(product_model.Sku)
	err := database.GetDB(constants.DatabaseConfigKey).First(sku, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return sku, nil
}

func FindByIds(ids []uint64) ([]*product_model.Sku, error) {
	skus := make([]*product_model.Sku, 0, len(ids))
	if len(ids) == 0 {
		return skus, nil
	}
	err := database.GetDB(constants.DatabaseConfigKey).Model(&product_model.Sku{}).Find(&skus, "id in (?)", ids).Error
	return skus, err
}

func DecrStock(spuId, skuId uint64, count uint64) (uint, error) {
	newStock, err := stock.DecrRealStock(spuId, skuId, count)
	if err == nil {
		asyncSetStock(skuId, spuId)
	}
	return newStock, err
}

func IncrStock(spuId, skuId uint64, count uint64) (uint, error) {
	newStock, err := stock.IncrRealStock(spuId, skuId, count)
	if err == nil {
		asyncSetStock(skuId, spuId)
	}
	return newStock, err
}

func asyncSetStock(skuId uint64, spuId uint64) {
	// 异步同步到db中的sku与spu，后续多sku时，db中取spu下sku库存之和同步给spu
	go func() {
		realStock, err := stock.GetRealStock(spuId, skuId)
		if err != nil {
			return
		}
		if err := database.GetDB(constants.DatabaseConfigKey).Model(&product_model.Sku{}).Where("id = ?", skuId).Update("stock", realStock).Error; err != nil {
			logrus.WithError(err).WithField("skuId", skuId).Warn("sync sku stock failed")
			return
		}
		if err := database.GetDB(constants.DatabaseConfigKey).Model(&product_model.Spu{}).Where("id = ?", spuId).Update("stock", realStock).Error; err != nil {
			logrus.WithError(err).WithField("spuId", spuId).Warn("sync spu stock failed")
			return
		}
	}()
}
