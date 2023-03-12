package config

import (
	"context"
	"encoding/json"

	"log"
	"sort"
	"time"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	fiberotel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	appConfig  *CommonConfig
	ParentCtx  context.Context
	ParentSpan trace.Span
)

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
	//tracing
	Tracesink           = "TRACING_SINK"
	DDAgentHost         = "DD_AGENT_HOST"
	DDAgentPort         = "DD_AGENT_PORT"
	JaegerAgentEndpoint = "JAEGER_AGENT_ENDPOINT"
	JaegerAgentHost     = "JAEGER_AGENT_HOST"
	JaegerAgentPort     = "JAGER_AGENT_PORT"
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
	Pprofenabled                 bool
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
	DatadogCofigHost            string
	DatadogConfigPort           string
	Tracesink                   string
	JaegerConfigEndpoint        string
	JeagerConfigHost            string
	JeagerConfigPort            string
	OpenTracingServiceName      string
	AppEnviroment               string
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

	//datadog config
	viper.SetDefault(Tracesink, "jaeger")
	viper.SetDefault(DDAgentHost, "127.0.0.1")
	viper.SetDefault(DDAgentPort, "8126")

	//jaeger config
	viper.SetDefault(JaegerAgentEndpoint, "http://jaeger:14268/api/traces")
	viper.SetDefault(JaegerAgentHost, "jaeger")
	viper.SetDefault(JaegerAgentPort, "6831")
}

func (config *CommonConfig) BuildConfig() *CommonConfig {

	config.setDefaultConfig()
	logConfig, logger := config.setLogConfig()

	config = &CommonConfig{
		LogConfig:                    logConfig,
		ChildFiberProcessIdleTimeout: viper.GetDuration(ChildFiberProcessIdleTimeout),
		SrvListenPort:                viper.GetString(SrvListenPort),
		Pprofenabled:                 viper.GetBool(Pprofenabled),
	}

	configJsonPresntation, _ := json.Marshal(config)
	logger.Info("Settup Config", zap.String("AppConfig", string(configJsonPresntation)))

	appConfig = config
	config.buildLogger()

	//initialize datadog
	if config.Tracesink == constant.DatadogTracingSink {
		appConfig.InitDatadogConfig()
	}

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
		Tracesink:                   viper.GetString(Tracesink),
		DatadogCofigHost:            viper.GetString(DDAgentHost),
		DatadogConfigPort:           viper.GetString(DDAgentPort),
		JaegerConfigEndpoint:        viper.GetString(JaegerAgentEndpoint),
		JeagerConfigHost:            viper.GetString(JaegerAgentHost),
		JeagerConfigPort:            viper.GetString(JaegerAgentPort),
		OpenTracingServiceName:      viper.GetString("OPEN_TRACING_SERVICE_NAME"),
		AppEnviroment:               viper.GetString("APP_ENVIRONMENT"),
	}

	return logConfig, configLogger
}

func (config CommonConfig) buildLogger() {

	ZapLogLevel := map[string]zapcore.Level{
		"DEBUG":  zapcore.DebugLevel,
		"INFO":   zapcore.InfoLevel,
		"WARN":   zapcore.WarnLevel,
		"ERROR":  zapcore.ErrorLevel,
		"FATAL":  zapcore.FatalLevel,
		"PANIC":  zapcore.PanicLevel,
		"DPANIC": zapcore.DPanicLevel,
	}

	var (
		logLevel  zapcore.Level = ZapLogLevel[config.LogLevel]
		err       error         = nil
		core      zapcore.Core
		zapLogger *zap.Logger
	)

	if logLevel == 0 {
		log.Fatalf("can't initialize zap logger - unsupported log level %v", logLevel)
	}

	if config.LogDestination == File {
		LogConfig := zap.NewDevelopmentEncoderConfig()
		LogConfig.FunctionKey = "F"

		var enc zapcore.Encoder

		if config.LogFormat == Json {
			enc = zapcore.NewJSONEncoder(LogConfig)
		} else {
			enc = zapcore.NewConsoleEncoder(LogConfig)
		}

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.LogFileName,
			MaxSize:    config.LogMaxSizeMb,
			MaxBackups: config.LogMaxBackupDays,
			MaxAge:     config.LogMaxAgeDaysBeforeRollover,
			Compress:   config.LogCompression,
		})

		core = zapcore.NewCore(enc, w, logLevel)
		zapLogger = zap.New(core)
	} else {
		logConfig := zap.NewDevelopmentConfig()
		logConfig.Level = zap.NewAtomicLevelAt(logLevel)
		logConfig.Encoding = config.LogFormat
		zapLogger, err = logConfig.Build()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger %v", err)
	}

	utils.Logger = zapLogger

}

func GetConfig() *CommonConfig {
	return appConfig
}

func InitConfig() {
	config := &CommonConfig{}
	appConfig = config.BuildConfig()
}

func (config *CommonConfig) InitDatadogConfig() {
	utils.Logger.Debug("Tracesink config", zap.String("commonConfig.Tra", config.Tracesink), zap.String("DatadogTracingSink", constant.DatadogTracingSink))
	agentAddress := config.DatadogCofigHost + ":" + config.DatadogConfigPort
	tracer.Start(tracer.WithAgentAddr(agentAddress))
	utils.Logger.Debug("TraceAgent", zap.String("agentAddress", agentAddress))
}

type OpentelemetryParantCtx struct {
	ParentCtx  context.Context
	ParentSpan trace.Span
}

func (parent *OpentelemetryParantCtx) SetOpentelementryParentCtx() {
	ParentCtx = parent.ParentCtx
	ParentSpan = parent.ParentSpan
}

func (config CommonConfig) setOpentelementry(app fiber.Router, ctx context.Context) (*tracesdk.TracerProvider, error) {

	var (
		tp      *tracesdk.TracerProvider
		sampler = tracesdk.AlwaysSample()
	)

	if config.Tracesink == constant.DatadogTracingSink {

		//datadog config
		agentAddress := config.DatadogCofigHost + ":" + config.DatadogConfigPort
		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(agentAddress))

		if err != nil {
			log.Panic(err.Error())
			return nil, err
		}

		tp = tracesdk.NewTracerProvider(
			//always be sure to batch in production
			tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(traceExporter)),
			// record information about this application in a resource
			tracesdk.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.OpenTracingServiceName),
				attribute.String("enviroment", config.AppEnviroment),
				//can add more attribute specific for ENvironment
			)),
			tracesdk.WithSampler(sampler),
		)

	} else {
		//jaeger config
		if len(config.JaegerConfigEndpoint) == 0 && (len(config.JeagerConfigHost) == 0 || len(config.JeagerConfigPort) == 0) {
			utils.Logger.Debug("properties are empty for config jeager-opentelementry will not config tracing")
			sampler = tracesdk.NeverSample()
		}

		var (
			exporters *jaeger.Exporter
			err       error
		)

		if len(config.JaegerConfigEndpoint) > 0 {
			utils.Logger.Debug("config jeager-opentelementr withcollectorEndpoint")
			exporters, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerConfigEndpoint)))
		} else {
			utils.Logger.Debug("config jeager-opentelementr withcollectorEndpoint")
			exporters, err = jaeger.New(jaeger.WithAgentEndpoint(
				jaeger.WithAgentPort(config.JeagerConfigPort),
				jaeger.WithAgentHost(config.JeagerConfigHost)))
		}

		if err != nil {
			log.Panic(err.Error())
			return nil, err
		}

		tp = tracesdk.NewTracerProvider(
			//always be sure to batch in product
			tracesdk.WithBatcher(exporters),
			// record info about this application in a resource
			tracesdk.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.OpenTracingServiceName),
				attribute.String("enviroment", config.AppEnviroment),
				// //can add more attribute specific for ENvironment
			)),
			tracesdk.WithSampler(sampler),
		)

	}

	//middleware use
	app.Use(fiberotel.New(fiberotel.Config{
		SpanName: "HTTP {{ .Method }} URL {{ .Path }}",
		Tracer:   otel.GetTracerProvider().Tracer("FiberComponet"),
		TracerStartAttributes: []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithNewRoot(),
		},
	}))

	otel.SetTracerProvider(tp)

	return tp, nil
}
