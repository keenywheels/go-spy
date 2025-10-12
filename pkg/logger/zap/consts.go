package zap

// consts for mode option
const (
	ProductionMode  = "production"
	DevelopmentMode = "development"
)

// consts for encoding option
const (
	ConsoleEncoding = "console"
	JsonEncoding    = "json"
)

// default path for logging
const defaultLogPath = "./log/app.log"

// default log params
const (
	defaultLogLvl     = "debug"
	defaultMaxLogSize = 10
	defaultMaxBackups = 5
	defaultMaxLogAge  = 30
)
