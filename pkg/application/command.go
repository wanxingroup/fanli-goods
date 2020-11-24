package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/spf13/viper"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/category"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/spu"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/state"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/utils/log"
)

func Start() {

	app := launcher.NewApplication(
		launcher.SetApplicationDescription(
			&launcher.ApplicationDescription{
				ShortDescription: "merchant service",
				LongDescription:  "support merchant user data management function. like merchant core data and chain store data",
			},
		),
		launcher.SetApplicationLogger(log.GetLogger()),
		launcher.SetApplicationEvents(
			launcher.NewApplicationEvents(
				launcher.SetOnInitEvent(func(app *launcher.Application) {

					unmarshalConfiguration()

					registerMerchantRPCRouter(app)

					idcreator.InitCreator(app.GetServiceId())
				}),
				launcher.SetOnStartEvent(func(app *launcher.Application) {

					migrateDatabases()
				}),
			),
		),
	)

	app.Launch()
}

func registerMerchantRPCRouter(app *launcher.Application) {

	rpcService := app.GetRPCService()
	if rpcService == nil {

		log.GetLogger().WithField("stage", "onInit").Error("get rpc service is nil")
		return
	}

	protos.RegisterSPUServer(rpcService.GetRPCConnection(), &spu.Controller{})
	protos.RegisterCategoryServer(rpcService.GetRPCConnection(), &category.Controller{})
	protos.RegisterStatusControllerServer(rpcService.GetRPCConnection(), &state.Controller{})
}

func unmarshalConfiguration() {
	err := viper.Unmarshal(config.Config)
	if err != nil {

		log.GetLogger().WithError(err).Error("unmarshal config error")
	}
}
