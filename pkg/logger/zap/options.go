package zap

type Option func(*Logger)

// LogLvl set loglvl
func LogLvl(loglvl string) Option {
	return func(l *Logger) {
		l.loglvl = loglvl
	}
}

// Mode set mode
func Mode(mode string) Option {
	return func(l *Logger) {
		l.mode = mode
	}
}

// Encoding set encoding
func Encoding(encoding string) Option {
	return func(l *Logger) {
		l.encoding = encoding
	}
}

// LogPath set logPath
func LogPath(path string) Option {
	return func(l *Logger) {
		l.logPath = path
	}
}

// MaxLogSize set maxLogSize
func MaxLogSize(size int) Option {
	return func(l *Logger) {
		l.maxLogSize = size
	}
}

// MaxLogBackups set maxLogBackups
func MaxLogBackups(size int) Option {
	return func(l *Logger) {
		l.maxLogBackups = size
	}
}

// MaxLogAge set maxLogAge
func MaxLogAge(size int) Option {
	return func(l *Logger) {
		l.maxLogAge = size
	}
}
