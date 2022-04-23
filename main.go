package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/nzai/tcst/api"
	"github.com/nzai/tcst/config"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	configPath    = flag.String("c", "config.toml", "toml config file path")
	logPath       = flag.String("log", "log.txt", "log file path")
	listenAddress = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()

	logger, err := initLogger(*logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// read config from file
	conf, err := config.Parse(*configPath)
	if err != nil {
		zap.L().Fatal("read config failed", zap.Error(err))
	}

	client := sts.NewClient(conf.SecretID, conf.SecretKey, nil)
	option := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          conf.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"*",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + conf.Region + ":uid/" + conf.AppID + ":" + conf.Bucket + "/*",
					},
				},
			},
		},
	}

	api.NewServer(*listenAddress, client, option).Run()

	zap.L().Info("http server stopped")
}

func initLogger(logPath string) (*zap.Logger, error) {
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30, // days
	})

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, infoPriority),
		zapcore.NewCore(fileEncoder, fileWriter, debugPriority),
	)

	return zap.New(core, zap.AddCaller()), nil
}
