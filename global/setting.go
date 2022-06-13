package global

import (
	"go_webapp/pkg/settings"
)

var (
	ServerSetting *settings.ServerSettingS
	AppSetting    *settings.AppSettingS
	LogSetting    *settings.LogSettingS
	MySqlSetting  *settings.MySqlSettingS
	RedisSetting  *settings.RedisSettingS
	JWTSetting    *settings.JWTSettingS
)
