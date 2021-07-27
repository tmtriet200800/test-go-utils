package pkgLogger

type LoggerInterface interface {
	Info(verbose bool, data interface{})
	Debug(verbose bool, data interface{})
	Error(verbose bool, data interface{})
}
