package zap

func (l *Logger) Debug(args ...interface{}) {
	l.sl.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sl.Debugf(template, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sl.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sl.Infof(template, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sl.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sl.Warnf(template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sl.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sl.Errorf(template, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sl.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sl.Panicf(template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sl.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sl.Fatalf(template, args...)
}
