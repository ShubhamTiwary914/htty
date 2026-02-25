package htty 

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	LOG_INFO  = "info"
	LOG_WARN  = "warn"
	LOG_ERROR = "error"
	LOG_DEBUG = "debug"
	LOG_ALL   = "all"
	LOG_ENVNAME = "LOGLEVEL"
)

var logger = log.New(os.Stdout, "", 0)
var devmodes = []string {LOG_INFO, LOG_WARN, LOG_ERROR, LOG_DEBUG, LOG_ALL}


func Logf(level string, format string, args ...interface{}) {
	if !assertAllowedLogLevel(level) {
		return; 
	}
	ts := time.Now().Format("02-01-2006(15:04:05.000)")
	//source = line that called this method
	pc, _, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			caller = fmt.Sprintf("%s:%d", fn.Name(), line)
		}
	}
	msg := fmt.Sprintf(format, args...)
	logger.Printf("level=%s ts=%s caller=%s msg=%q", level, ts, caller, msg)
}

func Debugf(format string, args ...interface{}) {
	Logf(LOG_DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(LOG_INFO, format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(LOG_ERROR, format, args...)
}


//redirect the logs to some debug file (suggested use only during debug MODE)
func RedirectLogs_toFile(outFile string, overwrite bool) *os.File {
	var flags int
	if overwrite {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}
	logfile, err := os.OpenFile(outFile, flags, 0644)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(logfile)
	return logfile	
}


//check if the log is allowed in this "LOGLEVEL" 
//(ex LOGLEVEL=all means allow all, LOGLEVEL=debug means allow only debug) 
func assertAllowedLogLevel(level string) bool {
	if os.Getenv(LOG_ENVNAME) == LOG_ALL {
		return true;	
	}		
	if os.Getenv(LOG_ENVNAME) == level {
		return true;	
	}
	return false
}

func pathBase(fp string) string {
	lastslash := -1
	for id := len(fp) - 1; id >= 0; id-- {
		if fp[id] == '/' {
			lastslash = id
			break
		}
	}
	if lastslash >= 0 && lastslash+1 < len(fp) {
		return fp[lastslash+1:]
	}
	return fp
}
