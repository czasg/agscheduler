package agscheduler

import "github.com/sirupsen/logrus"

var (
	AGSLog = logrus.New()
	Log    = AGSLog.WithFields(logrus.Fields{
		"AGSVersion": Version,
	})
)
