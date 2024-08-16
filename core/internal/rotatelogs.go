package internal

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"

	"server/global"
)

type lumberjackLogs struct{}

var LumberjackLogs = new(lumberjackLogs)

// GetWriteSyncer 获取 zapcore.WriteSyncer
func (l *lumberjackLogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := &lumberjack.Logger{
		Filename:   path.Join(global.Config.Zap.Director, level+".log"),
		MaxSize:    global.Config.RotateLogs.MaxSize,
		MaxBackups: global.Config.RotateLogs.MaxBackups,
		MaxAge:     global.Config.RotateLogs.MaxAge,
		Compress:   global.Config.RotateLogs.Compress,
	}

	if global.Config.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
