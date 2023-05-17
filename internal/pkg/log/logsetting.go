package Mylog

type LogSetting struct {
	RootPath     string
	SavePath     string
	SaveFilename string
	TimeFormat   string
	MaxSize      int32
	MaxBackups   int32
	Compress     bool
	JsonFormat   bool
	ShowLine     bool
	LogInConsole bool
	Level        string
	MaxDays      int32
}

type Option interface {
	apply(*LogSetting)
}

type optionFunc func(*LogSetting)

func (f optionFunc) apply(log *LogSetting) {
	f(log)
}

func NewLogSetting(options ...Option) *LogSetting {
	logSetting := &LogSetting{}
	for _, option := range options {
		option.apply(logSetting)
	}
	return logSetting
}

func WithLogRootPath(RootPath string) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.RootPath = RootPath
	})
}

func WithLogSavePath(SavePath string) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.SavePath = SavePath
	})
}

func WithLogSaveFilename(SaveFilename string) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.SaveFilename = SaveFilename
	})
}

func WithLogTimeFormat(TimeFormat string) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.TimeFormat = TimeFormat
	})
}

func WithLogMaxSize(MaxSize int32) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.MaxSize = MaxSize
	})
}

func WithLogMaxBackups(MaxBackups int32) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.MaxBackups = MaxBackups
	})
}

func WithLogCompress(Compress bool) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.Compress = Compress
	})
}

func WithLogJsonFormat(JsonFormat bool) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.JsonFormat = JsonFormat
	})
}

func WithLogShowLine(ShowLine bool) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.ShowLine = ShowLine
	})
}

func WithLogLogInConsole(logInConsole bool) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.LogInConsole = logInConsole
	})
}

func WithLogMaxDays(maxDays int32) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.MaxDays = maxDays
	})
}

func WithLogLevel(level string) Option {
	return optionFunc(func(setting *LogSetting) {
		setting.Level = level
	})
}
