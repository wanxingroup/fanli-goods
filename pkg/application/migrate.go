package application

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func migrateDatabases() {

	log.GetLogger().Info("starting to migrate database")
	migrateDatabaseAndLogError(product.Category{})
	migrateDatabaseAndLogError(product.Spu{})
	migrateDatabaseAndLogError(product.Sku{})
	log.GetLogger().Info("migrate database succeed")
}

func migrateDatabaseAndLogError(object interface{}) {

	if err := database.GetDB(constants.DatabaseConfigKey).AutoMigrate(object).Error; err != nil {

		log.GetLogger().WithField("object", fmt.Sprintf("%T", object)).
			WithError(err).Error("migrate failed")
	}
}
