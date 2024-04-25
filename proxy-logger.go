package main

import "fmt"

var ProxyLoggerLevelOverride map[string]LoggerLevel = map[string]LoggerLevel{}
var ProxyLoggerNames map[string]bool = map[string]bool{}

var ProxyLoggerLevelStringToLevel map[string]LoggerLevel = map[string]LoggerLevel{
	"0":     TraceLogLevel,
	"t":     TraceLogLevel,
	"tr":    TraceLogLevel,
	"tra":   TraceLogLevel,
	"trace": TraceLogLevel,

	"1":     DebugLogLevel,
	"d":     DebugLogLevel,
	"deb":   DebugLogLevel,
	"dbg":   DebugLogLevel,
	"debug": DebugLogLevel,

	"2":    InfoLogLevel,
	"i":    InfoLogLevel,
	"inf":  InfoLogLevel,
	"info": InfoLogLevel,

	"3":     PrintLogLevel,
	"p":     PrintLogLevel,
	"pr":    PrintLogLevel,
	"prn":   PrintLogLevel,
	"pri":   PrintLogLevel,
	"print": PrintLogLevel,

	"4":       WarningLogLevel,
	"w":       WarningLogLevel,
	"wrn":     WarningLogLevel,
	"warn":    WarningLogLevel,
	"warning": WarningLogLevel,

	"5":     ErrorLogLevel,
	"e":     ErrorLogLevel,
	"er":    ErrorLogLevel,
	"err":   ErrorLogLevel,
	"error": ErrorLogLevel,

	"6":     FatalLogLevel,
	"f":     FatalLogLevel,
	"fat":   FatalLogLevel,
	"ftl":   FatalLogLevel,
	"fatal": FatalLogLevel,

	"7":     PanicLogLevel,
	"pn":    PanicLogLevel,
	"pan":   PanicLogLevel,
	"panic": PanicLogLevel,
}

type ProxyBaseLogger LevelLoggerCompatible

type ProxyLoggerFlowPrefix struct {
	Source  string
	Target  string
	Reverse bool
}

type ProxyLoggerFlow = ProxyLoggerFlowPrefix

func (c *ProxyLoggerFlowPrefix) String() string {
	var s string
	if len(c.Source) > 0 {
		s = c.Source
		if len(c.Target) > 0 {
			var transfer string
			if c.Reverse {
				transfer = "<-"
			} else {
				transfer = "->"
			}
			s = s + transfer + c.Target // source --> target | source <-- target
		}
		// s = "[" + s + "]" // wrap inside []
	}
	// if len(c.Prefix) > 0 {
	// 	s = s + " " + c.Prefix
	// }
	return s
}

func (c *ProxyLoggerFlowPrefix) Set(source string, target string, reverse bool) *ProxyLoggerFlowPrefix {
	return &ProxyLoggerFlowPrefix{Source: source, Target: target, Reverse: reverse}
}

func (c *ProxyLoggerFlowPrefix) SetReverse(reverse bool) *ProxyLoggerFlowPrefix {
	return &ProxyLoggerFlowPrefix{Source: c.Source, Target: c.Target, Reverse: reverse}
}

func (c *ProxyLoggerFlowPrefix) SetSource(source string) *ProxyLoggerFlowPrefix {
	return &ProxyLoggerFlowPrefix{Source: source, Target: c.Target, Reverse: c.Reverse}
}

func (c *ProxyLoggerFlowPrefix) SetTarget(target string) *ProxyLoggerFlowPrefix {
	return &ProxyLoggerFlowPrefix{Source: c.Source, Target: target, Reverse: c.Reverse}
}

func (c *ProxyLoggerFlowPrefix) SetTargetDirection(target string, reverse bool) *ProxyLoggerFlowPrefix {
	return &ProxyLoggerFlowPrefix{Source: c.Source, Target: target, Reverse: reverse}
}

type LoggerCallStack struct {
	stack string
}

func (s *LoggerCallStack) String() string {
	return s.stack
}

func (s *LoggerCallStack) fork(frame string) *LoggerCallStack {
	if len(frame) == 0 {
		return s
	}
	var stack string
	if len(s.stack) > 0 {
		stack = s.stack + ">" + frame
	} else {
		stack = frame
	}
	return &LoggerCallStack{stack}
}

func (s *LoggerCallStack) Add(frame string) {
	if len(s.stack) > 0 && len(frame) > 0 {
		s.stack = s.stack + ":" + frame
	} else if len(frame) > 0 {
		s.stack = frame
	} else {
		s.stack = s.stack
	}
}

func (s *LoggerCallStack) append(frame string) {
	s.stack += frame
}

func (s *LoggerCallStack) forkInto(into *LoggerCallStack, frame string) {
	into.stack = s.stack
	into.Add(frame)
	// if len(s.stack) > 0 && len(frame) > 0 {
	// 	into.stack = s.stack + ":" + frame
	// } else if len(frame) > 0 {
	// 	into.stack = frame
	// } else {
	// 	into.stack = s.stack
	// }
}

type LoggerName struct {
	name string
}

func (n *LoggerName) String() string {
	return n.name
}

func (n *LoggerName) fork(name string) *LoggerName {
	if len(name) == 0 {
		return n
	}
	var fullName string
	if len(n.name) > 0 {
		fullName = n.name + "." + name
		ProxyLoggerNames[fullName] = true
	} else {
		fullName = name
	}
	return &LoggerName{fullName}
}

func (n *LoggerName) forkInto(into *LoggerName, name string) {
	if len(n.name) > 0 && len(name) > 0 {
		into.name = n.name + "." + name
		ProxyLoggerNames[into.name] = true
	} else if len(name) > 0 {
		into.name = name
		ProxyLoggerNames[into.name] = true
	} else {
		into.name = n.name
	}
}

type errorFmt struct {
	*ProxyLogger
}

type ProxyLogger struct {
	ProxyBaseLogger
	logger LevelCoreLogger

	flowPrefix   ProxyLoggerFlowPrefix
	parentPrefix string
	prefix       string
	// _fullPrefix  string
	// _prefixBuilt bool
	composedPrefix string
	E              *errorFmt

	flow  ProxyLoggerFlowPrefix
	Name  LoggerName
	Stack LoggerCallStack
}

// this is just for setting options when creating new
type ProxyLoggerInfo struct {
	Level  LoggerLevel
	Flow   *ProxyLoggerFlow
	Name   string
	Frame  string
	Prefix string
}

func NewProxyLogger(logger ProxyBaseLogger) *ProxyLogger {
	// type LoggerLevelWrapper struct {
	// 	LevelCoreLogger
	// }
	loggerWrapper := &LoggerLevelWrapper{logger}
	proxyLogger := &ProxyLogger{loggerWrapper, logger, ProxyLoggerFlowPrefix{}, "", "", "", &errorFmt{}, ProxyLoggerFlowPrefix{}, LoggerName{}, LoggerCallStack{}}
	loggerWrapper.LevelCoreLogger = proxyLogger
	proxyLogger.SetLevel(logger.GetLevel())
	// proxyLogger._initErrorFmt()
	proxyLogger.E.ProxyLogger = proxyLogger
	return proxyLogger
	// return (&ProxyLogger{logger, ProxyLoggerFlowPrefix{}, "", "", "", false, &errorFmt{}}).SetLevel(logger.GetLevel())._initErrorFmt()
}

// func NewProxyLogger(logger ProxyBaseLogger, flowPrefix ProxyLoggerFlowPrefix, prefix string, level LoggerLevel) *ProxyLogger {
// 	return (&ProxyLogger{logger, flowPrefix, "", prefix, ""}).Build().SetLevel(logger.GetLevel())
// }

// func (l *ProxyLogger) NewLogger() *ProxyLogger {
// 	return (&ProxyLogger{l.ProxyBaseLogger, l.flowPrefix, l.prefix, "", ""})
// }

func (l *ProxyLogger) Flow() *ProxyLoggerFlowPrefix {
	return &l.flow
}

func (l *ProxyLogger) Fork(info ProxyLoggerInfo) *ProxyLogger {
	// name, frame, level, flow settings
	// flow settings: source, target, direction
	level := info.Level
	if level == UnsetLogLevel {
		level = l.logger.GetLevel()
	}
	logger := &LoggerLevelCoreWrapper{l.logger, level}
	loggerWrapper := &LoggerLevelWrapper{logger}
	proxyLogger := &ProxyLogger{loggerWrapper, logger, l.flowPrefix, l.prefix, l.prefix, "", &errorFmt{}, l.flow, LoggerName{}, LoggerCallStack{}}
	l.Name.forkInto(&proxyLogger.Name, info.Name)
	l.Stack.forkInto(&proxyLogger.Stack, info.Frame)
	if info.Flow != nil {
		proxyLogger.flow = *info.Flow
	}
	if len(info.Prefix) > 0 {
		proxyLogger.prefix = info.Prefix
	}
	loggerWrapper.SetWrappedLogger(proxyLogger)
	// no need for wildcards, because the level applies to all inherited loggers and any inherited logger can be also overrided
	if level, present := ProxyLoggerLevelOverride[proxyLogger.Name.String()]; present {
		loggerWrapper.SetLevel(level)
	}
	// loggerWrapper.SetLevel(l.GetLevel())
	// proxyLogger._initErrorFmt()
	proxyLogger.E.ProxyLogger = proxyLogger
	return proxyLogger
}

func (l *ProxyLogger) Derive() *ProxyLogger {
	return l.Fork(ProxyLoggerInfo{})
	// create new logger with the same level and base logger
	// log := &LoggerLevelWrapper{NewLoggerLevelCoreWrapper(l.logger)}
	logger := &LoggerLevelCoreWrapper{l.logger, l.logger.GetLevel()}
	loggerWrapper := &LoggerLevelWrapper{logger}
	// log := &LoggerLevelWrapper{NewLoggerLevelCoreWrapper(l.ProxyBaseLogger.GetLogger())}
	proxyLogger := &ProxyLogger{loggerWrapper, logger, l.flowPrefix, l.prefix, "", "", &errorFmt{}, ProxyLoggerFlowPrefix{}, LoggerName{}, LoggerCallStack{}}
	// log.LevelCoreLogger = proxyLogger
	loggerWrapper.SetWrappedLogger(proxyLogger)
	loggerWrapper.SetLevel(l.GetLevel())
	// proxyLogger._initErrorFmt()
	proxyLogger.E.ProxyLogger = proxyLogger
	return proxyLogger
	// return (&ProxyLogger{log, log, l.flowPrefix, l.prefix, "", "", false, &errorFmt{}})._initErrorFmt()
	// return (&ProxyLogger{log, l.flowPrefix, l.prefix, "", "", false, &errorFmt{}})._initErrorFmt()
}

func (l *ProxyLogger) _initErrorFmt() *ProxyLogger {
	l.E.ProxyLogger = l
	return l
}

func (l *ProxyLogger) WithExtension(prefixExtension string) *ProxyLogger {
	// proxyLogger := l.Fork(ProxyLoggerInfo{Prefix: l.prefix + prefixExtension})
	proxyLogger := l.Fork(ProxyLoggerInfo{})
	proxyLogger.Stack.append(prefixExtension)
	return proxyLogger
	return (&ProxyLogger{l.ProxyBaseLogger, l.logger, l.flowPrefix, l.parentPrefix, l.prefix + prefixExtension, "", &errorFmt{}, ProxyLoggerFlowPrefix{}, LoggerName{}, LoggerCallStack{}})._initErrorFmt()
	// return (&ProxyLogger{l.ProxyBaseLogger, l.flowPrefix, l.parentPrefix, l.prefix + prefixExtension, "", false, &errorFmt{}})._initErrorFmt()
}

//	func (l *ProxyLogger) SetLevel(level LoggerLevel) *ProxyLogger {
//		l.ProxyBaseLogger.SetLevel(level)
//		return l
//	}
func (l *ProxyLogger) SetLevel(level LoggerLevel) {
	// l.ProxyBaseLogger.SetLevel(level)
	l.logger.SetLevel(level)
}

func (l *ProxyLogger) GetLevel() LoggerLevel {
	return l.logger.GetLevel()
}

func (l *ProxyLogger) SetPrefix(prefix string) *ProxyLogger {
	// l.prefix = prefix
	l.Stack.Add(prefix)
	l.composedPrefix = ""
	// l._prefixBuilt = false
	return l
}

func (l *ProxyLogger) SetSource(source string) *ProxyLogger {
	l.flow.Source = source
	l.flow.Target = ""
	l.composedPrefix = ""
	// l._prefixBuilt = false
	return l
}

func (l *ProxyLogger) SetTarget(target string, reverse bool) *ProxyLogger {
	l.flow.Target = target
	l.flow.Reverse = reverse
	l.composedPrefix = ""
	// l._prefixBuilt = false
	return l
}

// func (l *ProxyLogger) Build() *ProxyLogger {
// 	p := l.flowPrefix.String()
// 	if len(l.parentPrefix) > 0 {
// 		if len(p) > 0 {
// 			p = p + " "
// 		}
// 		p = p + l.parentPrefix + ":"
// 	}
// 	if len(l.prefix) > 0 {
// 		if len(p) > 0 {
// 			p = p + " "
// 		}
// 		p = p + l.prefix + ":"
// 	}
// 	l.composedPrefix = p
// 	return l
// }

func (l *ProxyLogger) getPrefix() string {
	// if !l._prefixBuilt {
	if len(l.composedPrefix) == 0 {
		// l.composedPrefix = fmt.Sprintf("[%s] [%s]", l.flowPrefix.String(), l.Stack.String())
		// p := l.flowPrefix.String()
		// if len(l.parentPrefix) > 0 {
		// 	if len(p) > 0 {
		// 		p = p + " "
		// 	}
		// 	p = p + l.parentPrefix + ":"
		// }
		// if len(l.prefix) > 0 {
		// 	if len(p) > 0 {
		// 		p = p + " "
		// 	}
		// 	p = p + l.prefix + ":"
		// }
		// l.composedPrefix = p
		// l._prefixBuilt = true
		// l.composedPrefix = p + fmt.Sprintf("| [%s] [%s]", l.flow.String(), l.Stack.String())
		flow := l.flow.String()
		stack := l.Stack.String()
		if len(flow) > 0 && len(stack) > 0 {
			l.composedPrefix = fmt.Sprintf("[%s] [%s]", flow, stack)
		} else if len(flow) > 0 {
			l.composedPrefix = fmt.Sprintf("[%s]", flow)
		} else if len(stack) > 0 {
			l.composedPrefix = fmt.Sprintf("[%s]", stack)
		} else {
			l.composedPrefix = " "
		}
		// l.composedPrefix = fmt.Sprintf("[%s] [%s]", l.flow.String(), l.Stack.String())
		// l.composedPrefix = fmt.Sprintf("[%s] [%s] [%s]", l.Name.String(), l.flow.String(), l.Stack.String())
		if len(l.prefix) > 0 {
			l.composedPrefix += " " + l.prefix
		}
	}
	return l.composedPrefix
}

func (l *ProxyLogger) Prefix() string {
	return l.prefix
}

func (l *ProxyLogger) LLog(level LoggerLevel, v ...any) {
	v = append([]any{l.getPrefix()}, v...)
	l.logger.LLog(level, v...)
}

func (l *ProxyLogger) LLogf(level LoggerLevel, format string, v ...any) {
	l.logger.LLogf(level, l.getPrefix()+" "+format, v...)
}

/*
func (l *ProxyLogger) Trace(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Trace(v...)
}

func (l *ProxyLogger) Tracef(format string, v ...any) {
	l.ProxyBaseLogger.Tracef(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Debug(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Debug(v...)
}

func (l *ProxyLogger) Debugf(format string, v ...any) {
	l.ProxyBaseLogger.Debugf(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Info(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Info(v...)
}

func (l *ProxyLogger) Infof(format string, v ...any) {
	l.ProxyBaseLogger.Infof(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Warn(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Warn(v...)
}

func (l *ProxyLogger) Warnf(format string, v ...any) {
	l.ProxyBaseLogger.Warnf(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Error(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Error(v...)
}

func (l *ProxyLogger) Errorf(format string, v ...any) {
	l.ProxyBaseLogger.Errorf(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Fatal(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Fatal(v...)
}

func (l *ProxyLogger) Fatalf(format string, v ...any) {
	l.ProxyBaseLogger.Fatalf(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Panic(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Panic(v...)
}

func (l *ProxyLogger) Panicf(format string, v ...any) {
	l.ProxyBaseLogger.Panicf(l.fullPrefix()+" "+format, v...)
}

// compatilbity
func (l *ProxyLogger) Print(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Print(v...)
}

func (l *ProxyLogger) Printf(format string, v ...any) {
	l.ProxyBaseLogger.Printf(l.fullPrefix()+" "+format, v...)
}

func (l *ProxyLogger) Println(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Println(v...)
}

func (l *ProxyLogger) Fatalln(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Fatalln(v...)
}

func (l *ProxyLogger) Panicln(v ...any) {
	v = append([]any{l.fullPrefix()}, v...)
	l.ProxyBaseLogger.Panicln(v...)
}
*/

func (e *errorFmt) Errorf(format string, v ...any) error {
	return fmt.Errorf(e.ProxyLogger.getPrefix()+": "+format, v...)
	// return fmt.Errorf(e.ProxyLogger.Prefix()+": "+format, v...)
}

// passthrough
func (e *errorFmt) Sprintf(format string, v ...any) string {
	return fmt.Sprintf(format, v...)
}
