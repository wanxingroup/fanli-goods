package test

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
)

func Release() {

	var err error
	err = database.GetDB(constants.DatabaseConfigKey).Exec(fmt.Sprintf("DROP DATABASE %s", databaseName)).Error
	if err != nil {
		logrus.WithField("error", err).Error("drop database error")
	}

	err = database.Disconnect(constants.DatabaseConfigKey)
	if err != nil {
		logrus.WithField("error", err).Error("close database error")
	}
}
