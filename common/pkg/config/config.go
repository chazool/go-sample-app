package config

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var appConfig *CommonConfig

// configuration constance
const (
	SrvListenPort                = "SRV_LISTEN_PORT"
	ChildFiberProcessIdleTimeout = "CHILD_FIBER_PROCESS_IDLE_TIMEOUT"
	LogDestination               = "LOG_DESTINATION"
	//log constance
	LogFileName                 = "LOG_FILE_NAME"
	LogMaxSizeMb                = "LOG_MAX_SIZE_MB"
	LogMaxBackupDays            = "LOG_MAX_BACKUP_DAYS"
	LogMaxAgeDaysBeforeRollover = "LOG_MAX__AGE_DAYS_BEFORE_ROLLOVER"
	LogCompressionEnabled       = "LOG_COMPOSITION_ENABLED"
	LogLevel                    = "LOG_LEVEL"
	LogFormat                   = "LOG_FORMAT"
	Pprofenabled                = "PPRO_ENABLED"
	//log constance values
	Console = "console"
	File    = "file"
	Debug   = "DEBUG"
	Json    = "json"
)

type CommonConfig struct {
	_ struct{}
	LogConfig
	SrvListenPort                string
	ChildFiberProcessIdleTimeout time.Duration
}

type LogConfig struct {
	_                           struct{}
	LogDestination              string
	LogFileName                 string
	LogMaxSizeMb                int
	LogMaxBackupDays            int
	LogMaxAgeDaysBeforeRollover int
	LogCompression              bool
	LogLevel                    string
	LogFormat                   string
}

// setDefaultConfig is using added application default configurations
func (config *CommonConfig) setDefaultConfig() {

	viper.SetDefault(SrvListenPort, "8090")

	dur, _ := time.ParseDuration("10s")
	viper.SetDefault(ChildFiberProcessIdleTimeout, dur)

	//log default config
	//you can suppy "console" or "File". if console, logging goes to stdout, if tile, goes to LOG_FILE_NAME
	viper.SetDefault(LogDestination, Console)
	viper.SetDefault(LogFileName, "app.log")
	viper.SetDefault(LogMaxSizeMb, 100)
	viper.SetDefault(LogMaxBackupDays, 30)
	viper.SetDefault(LogMaxAgeDaysBeforeRollover, 28)
	viper.SetDefault(LogCompressionEnabled, true)
	viper.SetDefault(LogLevel, Debug)
	viper.SetDefault(LogFormat, Console)

	//you can suppy "console" or "File". if json, loggin formant is in json
	viper.SetDefault(LogFileName, Json)
	viper.SetDefault(Pprofenabled, "true")

}

func (config *CommonConfig) BuildConfig() *CommonConfig {

	config.setDefaultConfig()
	logConfig, logger := config.setLogConfig()

	config = &CommonConfig{
		LogConfig:                    logConfig,
		ChildFiberProcessIdleTimeout: viper.GetDuration(ChildFiberProcessIdleTimeout),
		SrvListenPort:                viper.GetString(SrvListenPort),
	}

	configJsonPresntation, _ := json.Marshal(config)
	logger.Info("Settup Config", zap.String("AppConfig", string(configJsonPresntation)))

	appConfig = config
	return config
}

// setLogConfig is using settup the zap logger configuration
func (config *CommonConfig) setLogConfig() (LogConfig, *zap.Logger) {

	configLogger, _ := zap.NewDevelopmentConfig().Build()

	defer configLogger.Sync()

	supportedLogDestinations := []string{Console, File}
	sort.Strings(supportedLogDestinations)
	logDestination := viper.GetString(LogDestination)

	if !utils.Contains(supportedLogDestinations, logDestination) {
		configLogger.Fatal("Invalid log destination specified", zap.String(LogDestination, logDestination), zap.Strings("supportedLogDestinations", supportedLogDestinations))
	} else {
		configLogger.Info("log destination is set to ", zap.String(LogDestination, logDestination))
	}

	supportedLogFormats := []string{Console, Json}
	logFormat := viper.GetString(LogFormat)

	if !utils.Contains(supportedLogFormats, logFormat) {
		configLogger.Fatal("Invalid Log Format specified", zap.String(LogFormat, logFormat), zap.Strings("supportedLogDestinations", supportedLogDestinations))
	} else {
		configLogger.Info("log format is se to", zap.String(LogFormat, logFormat))
	}

	logConfig := LogConfig{
		LogDestination:              logDestination,
		LogFormat:                   logFormat,
		LogFileName:                 viper.GetString(LogFileName),
		LogLevel:                    viper.GetString(LogLevel),
		LogMaxSizeMb:                viper.GetInt(LogMaxSizeMb),
		LogMaxBackupDays:            viper.GetInt(LogMaxBackupDays),
		LogMaxAgeDaysBeforeRollover: viper.GetInt(LogMaxAgeDaysBeforeRollover),
		LogCompression:              viper.GetBool(LogCompressionEnabled),
	}

	return logConfig, configLogger
}

func GetConfig() *CommonConfig {
	return appConfig
}

func InitConfig() {
	config := &CommonConfig{}
	appConfig = config.BuildConfig()
}
