package common

import (
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/odacremolbap/rest-demo/pkg/log"
	pkglogrus "github.com/odacremolbap/rest-demo/pkg/log/logrus"
)

// ConfigureLogrusLogger configures the logger to use throughout the application
func ConfigureLogrusLogger(format string, verbosity int) {

	switch strings.ToLower(format) {
	case "server":
		logrus.SetFormatter(pkglogrus.NewServerFomatter(false, ""))
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		logrus.Fatal("incorrect log formatter")
	}

	logger := pkglogrus.NewLogrusLogr(logrus.StandardLogger(), "", verbosity)
	log.SetDefaultLogger(logger)
}
