package initializers

import (
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
)

const requestLogFormat = " %s \"%s %s %s\" %d"

type logger struct {
	fileLogger *log.Logger
	sysLogger  *syslog.Writer
}

func (l *logger) Info(m string) {
	if l.sysLogger != nil {
		l.sysLogger.Info(m)
	} else if l.fileLogger != nil {
		l.fileLogger.Println(m)
	}
}

func (l *logger) Err(m string) {
	if l.sysLogger != nil {
		l.sysLogger.Err(m)
	} else if l.fileLogger != nil {
		l.fileLogger.Println(m)
	}
}

func (l *logger) Warning(m string) {
	if l.sysLogger != nil {
		l.sysLogger.Warning(m)
	} else if l.fileLogger != nil {
		l.fileLogger.Println(m)
	}
}

func NewLogger(fileLogger *log.Logger, syslog *syslog.Writer) *logger {
	return &logger{fileLogger, syslog}
}

func LogRequest(r *http.Request, status int, err error) {
	logMsg := fmt.Sprintf(requestLogFormat, r.RemoteAddr, r.Method, r.URL.String(), r.Proto, status)

	if err != nil {
		logMsg += fmt.Sprintf("\nQUERY PARAMS: %s", r.URL.RawQuery)

		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				logMsg += fmt.Sprintf("\nREQUEST BODY:\n%8s", body)
			}
		}

		logMsg += fmt.Sprintf("\nERROR: %s", err.Error())
	}

	Logger.Info(logMsg)
}
