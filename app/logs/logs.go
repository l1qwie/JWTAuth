package logs

import "log"

var (
	caption string
	DEBUG   bool
	INFO    bool
)

func SetDebug() {
	DEBUG = true
	caption = "DEBUG"
}

func StartPoint(point, method string) {
	log.Printf("[INFO] %s (%s) has been launched", point, method)
}

func ParameterIsRequired(param string) {
	if DEBUG {
		log.Printf("[%s] %s is required", caption, param)
	}
}
