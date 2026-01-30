package execute

import "log"

var (
	logger Logger = debugLogger{}
)

type Logger interface {
	Infof(string, ...interface{})
}

type debugLogger struct{}

func (debugLogger) Infof(format string, v ...interface{}) {
	if !debug {
		return
	}

	log.Printf(format, v...)
}
