package global

import (
	"blog_service/pkg/logger"
	"blog_service/pkg/setting"

	"github.com/jinzhu/gorm"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DataBaseSetting *setting.DataBaseSettingS
)

var (
	DBEngine *gorm.DB
	Logger   *logger.Logger
)
