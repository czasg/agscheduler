package agscheduler

var (
	Version = "0.0.1"
	Author  = "https://github.com/CzaOrz"
	GitHub  = "https://github.com/CzaOrz/agscheduler"
	GitBook = "https://github.com/CzaOrz/agscheduler/blob/main/README.md"
	Email   = "972542655@qq.com"
)

func init() {
	Log.WithFields(GenAGSDetails()).Debugln("Welcome to AGS ~")
}
