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

func IsntValid(param string) {
	if DEBUG {
		log.Printf("[%s] %s is not valid", caption, param)
	}
}

func HasDone(param string) {
	if DEBUG {
		log.Printf("[%s] %s has been done", caption, param)
	}
}

func HasntDone(param string) {
	if DEBUG {
		log.Printf("[%s] %s hasn't been done", caption, param)
	}
}

func UnknownParam(param string) {
	if DEBUG {
		log.Printf("[%s] %s is unkown", caption, param)
	}
}
