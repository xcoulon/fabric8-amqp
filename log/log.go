package log

import (
	"github.com/fabric8-services/fabric8-amqp/configuration"
	"github.com/sirupsen/logrus"
)

var config configuration.Config

func init() {
	config = configuration.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		// DisableSorting: true,
	})
}

// Infof displays the given msg with optional args at the `info` level,
// preceeded by the name of the pod in which the program is running.
func Infof(msg string, args ...interface{}) {
	logrus.WithField("id", config.GetPodName()).Infof(msg, args...)
}

// Warn displays the given msg at the `warn` level,
// preceeded by the name of the pod in which the program is running.
func Warn(msg string) {
	logrus.WithField("id", config.GetPodName()).Warn(msg)
}

// Warnf displays the given msg with optional args at the `warn` level,
// preceeded by the name of the pod in which the program is running.
func Warnf(msg string, args ...interface{}) {
	logrus.WithField("id", config.GetPodName()).Warnf(msg, args...)
}

// Errorf displays the given msg with optional args at the `warn` level,
// preceeded by the name of the pod in which the program is running.
func Errorf(msg string, args ...interface{}) {
	logrus.WithField("id", config.GetPodName()).Errorf(msg, args...)
}

// Fatal displays the given err at the `fatal` level,
// preceeded by the name of the pod in which the program is running.
func Fatal(err error) {
	logrus.WithField("id", config.GetPodName()).Fatal(err.Error())
}

// Fatalf displays the given err at the `fatal` level,
// preceeded by the name of the pod in which the program is running.
func Fatalf(msg string, args ...interface{}) {
	logrus.WithField("id", config.GetPodName()).Fatalf(msg, args...)
}
