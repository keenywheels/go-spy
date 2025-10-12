package zap

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wrapper over zap logger
type Logger struct {
	sl *zap.SugaredLogger

	loglvl   string
	mode     string // production/development
	encoding string // console/json

	logPath       string // path for log file, file's dir will be used as dir for other logs
	maxLogSize    int    // max size of log file in MB
	maxLogBackups int    // max number of log file backups
	maxLogAge     int    // max lifetime for log file in days
}

// New create new Logger instance
func New(opts ...Option) *Logger {
	l := &Logger{
		mode:          ProductionMode,
		encoding:      JsonEncoding,
		loglvl:        defaultLogLvl,
		logPath:       defaultLogPath,
		maxLogSize:    defaultMaxLogSize,
		maxLogBackups: defaultMaxBackups,
		maxLogAge:     defaultMaxLogAge,
	}

	for _, opt := range opts {
		opt(l)
	}

	// setup encoder config
	var encoderCfg zapcore.EncoderConfig
	if l.mode == DevelopmentMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// create encoder
	var encoder zapcore.Encoder

	if l.encoding == ConsoleEncoding {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	// create file writeSyncer
	fileWriterSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.logPath,
		MaxSize:    l.maxLogSize,
		MaxBackups: l.maxLogBackups,
		MaxAge:     l.maxLogAge,
		Compress:   false,
	})

	// create sugared zaplogger
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(fileWriterSyncer, zapcore.AddSync(os.Stdout)),
		zap.NewAtomicLevelAt(l.getLogLvl()),
	)

	l.sl = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	return l
}

// Sync wrapper over zaplogger Sync
func (l *Logger) Sync() error {
	return l.sl.Sync()
}

// ToStdLog return std *log.Logger
func (l *Logger) ToStdLog() *log.Logger {
	// return log.New if nil
	if l == nil || l.sl == nil {
		return log.New(nil, "", 0)
	}

	return zap.NewStdLog(l.sl.Desugar())
}
