package logger

type LogLevel uint8

const (
	Emergency LogLevel = 1 << iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug
)

func LogLevelDebug() LogLevel {
	return Debug | Info | Notice | Warning | Error | Critical | Alert | Emergency
}
func LogLevelInfo() LogLevel      { return Info | Notice | Warning | Error | Critical | Alert | Emergency }
func LogLevelNotice() LogLevel    { return Notice | Warning | Error | Critical | Alert | Emergency }
func LogLevelWarning() LogLevel   { return Warning | Error | Critical | Alert | Emergency }
func LogLevelError() LogLevel     { return Error | Critical | Alert | Emergency }
func LogLevelCritical() LogLevel  { return Critical | Alert | Emergency }
func LogLevelAlert() LogLevel     { return Alert | Emergency }
func LogLevelEmergency() LogLevel { return Emergency }

func (l LogLevel) Match(v LogLevel) bool {
	if l.Has(v) {
		return true
	}
	return false
}

func (l LogLevel) Has(v LogLevel) bool {
	return l == (l & v)
}

func (l LogLevel) String() string {
	var str string
	for i := 1; i <= 255; i <<= 1 {
		if l.Has(LogLevel(i)) {
			switch LogLevel(i) {
			case Emergency:
				str += "EMERGENCY|"
			case Alert:
				str += "ALERT|"
			case Critical:
				str += "CRITICAL|"
			case Error:
				str += "ERROR|"
			case Warning:
				str += "WARNING|"
			case Notice:
				str += "NOTICE|"
			case Info:
				str += "INFO|"
			case Debug:
				str += "DEBUG|"
			}
		}
	}
	if s := len(str); s > 0 {
		return str[:s-1]
	}
	return "UNKNOWN"
}
