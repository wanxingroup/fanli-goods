package config

import (
	launcherConfig "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/config"
)

type configuration struct {
	launcherConfig.StandardConfig `mapstructure:",squash"`
}

var (
	Config = new(configuration)
)
