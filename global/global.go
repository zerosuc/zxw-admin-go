package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"server/config"
)

var (
	Viper              *viper.Viper
	Config             config.Server
	Log                *zap.Logger
	Db                 *gorm.DB
	Redis              *redis.Client
	ConcurrencyControl = &singleflight.Group{}
	Cron               *cron.Cron
)
