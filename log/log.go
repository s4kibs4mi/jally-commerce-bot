package log

import (
	"github.com/s4kibs4mi/twilfe/log/hooks"
	"github.com/sirupsen/logrus"
	"os"
)

var defLogger *logrus.Logger

func init() {
	defLogger = logrus.New()
	defLogger.Out = os.Stdout
	defLogger.AddHook(hooks.NewHook())
}

func Log() *logrus.Logger {
	return defLogger
}
