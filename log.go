package agscheduler

import "github.com/sirupsen/logrus"

var (
	AGSLog = logrus.New()
	Log    = AGSLog.WithFields(GenAGSVersion())
)

func GenAGSVersion() logrus.Fields {
	return logrus.Fields{"AGSVersion": Version}
}

func GenASGModule(module string) logrus.Fields {
	return logrus.Fields{"ASGModule": module}
}

func GenAGSDetails() logrus.Fields {
	return logrus.Fields{
		"ASGAuthor":  Author,
		"AGSGitHub":  GitHub,
		"AGSGitBook": GitBook,
		"AGSEmail":   Email,
	}
}
