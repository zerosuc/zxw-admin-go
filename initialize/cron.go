package initialize

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"server/config"
	"server/global"
	modelSysTool "server/model/sysTool"
	"server/utils"
)

// Crontab 添加计划任务
func crontab() {
	if global.Config.Crontab.Open {
		ct := cron.New()
		for index := range global.Config.Crontab.Objects {
			go func(cObject config.Object) {
				_, err := ct.AddFunc(global.Config.Crontab.Spec, func() {
					err := utils.ClearTable(global.Db, cObject.TableName, cObject.CompareField, cObject.Interval)
					if err != nil {
						global.Log.Error("clear table", zap.Error(err))
					}
				})
				if err != nil {
					global.Log.Error("cron add func", zap.Error(err))
				}
			}(global.Config.Crontab.Objects[index])
		}
		// 启动cron
		ct.Start()
	}
}

// InitCron 初始化Cron
func InitCron() *cron.Cron {
	// 配置文件方式cron
	crontab()
	// 页面方式配置
	instance := cron.New(cron.WithSeconds()) // 支持秒
	instance.Start()                         // 启动cron
	return instance
}

func CheckCron() {
	var cronModelList []modelSysTool.CronModel
	global.Db.Where("open = ?", 1).Find(&cronModelList)
	for _, cronModel := range cronModelList {
		entryId, err := global.Cron.AddJob(cronModel.Expression, &cronModel)
		if err != nil {
			global.Log.Error("CRON", zap.Error(err))
		} else {
			global.Db.Model(cronModel).Update("entryId", entryId)
		}
	}
	global.Log.Info("CheckCron success")

}
